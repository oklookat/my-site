package article

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"servus/apps/elven/base"
	"servus/apps/elven/file"
	"servus/core/external/utils"
	"strconv"
	"strings"
)

// validate params to get paginated articles.
func ValidateGetParams(a *GetParams, params url.Values, isAdmin bool) {
	var err error

	// "published" param
	var published = params.Get("drafts")
	if len(published) < 1 {
		a.Drafts = false
	} else {
		a.Drafts, err = strconv.ParseBool(published)
		if err != nil {
			a.Drafts = false
		}
	}
	if a.Drafts && !isAdmin {
		a.Drafts = false
	}

	// "by" param.
	var by = params.Get("by")
	if by == "created" || by == "updated" {
		if !isAdmin {
			by = "published"
		}
	} else {
		by = "published"
	}
	by += "_at"
	a.By = by

	// "newest" param.
	var newest = params.Get("newest")
	if len(newest) < 1 {
		a.Newest = true
	} else {
		a.Newest, err = strconv.ParseBool(newest)
		if err != nil {
			a.Newest = true
		}
	}

	// "page" param
	var pageStr = params.Get("page")
	if len(pageStr) < 1 {
		a.Page = 1
	} else {
		var page = 0
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
		a.Page = page
	}

	// "title" param
	var title = params.Get("title")
	if len(title) > 0 {
		a.Title = &title
	}
}

// validate/filter body to change/create article.
//
// requestMethod = validation depends on request method (changing mode).
//
// body = request body.
//
// reference = reference article. Needs when we change existing article.
//
// ----returns:----
//
// nil if validation error;
//
// if changing mode = copy of reference but filtered/validated by body.
//
// if creating mode = semi-filled filtered/validated article by body.
func ValidateBody(requestMethod string, body io.ReadCloser, reference *Model) (filtered *Model, err error) {

	var validationErr = base.ValidationError{}

	// updating existing article?
	var isChangingMode = requestMethod == http.MethodPut || requestMethod == http.MethodPatch

	// changing mode and no reference? we not validate this.
	if isChangingMode && reference == nil {
		validationErr.New("body")("changing mode enabled but reference has nil pointer")
		err = &validationErr
		return
	}

	// decode body.
	var bodyStruct = &Body{}
	err = json.NewDecoder(body).Decode(bodyStruct)
	if err != nil {
		return
	}

	// FUNC title.
	//
	// returns: is valid, title length.
	var checkTitle = func() (bool, int) {
		var isNil = bodyStruct.Title == nil
		if isNil {
			return false, 0
		}
		var titleLength = utils.LenRune(*bodyStruct.Title)
		var isEmpty = titleLength < 1
		return !isEmpty, titleLength
	}

	// body has title field?
	var isTitle, titleLength = checkTitle()

	// pre check title.
	if isTitle && titleLength > 124 {
		validationErr.New("title")("max length is 124")
		err = &validationErr
		return
	}

	// FUNC isPublished.
	var checkIsPublished = func() bool {
		return bodyStruct.IsPublished != nil
	}

	// FUNC cover id.
	var checkCoverID = func() bool {
		return bodyStruct.CoverID != nil
	}

	// FUNC content.
	//
	// returns: is valid, content length.
	var checkContent = func() (bool, int) {
		var isNil = bodyStruct.Content == nil
		if isNil {
			return false, 0
		}
		var contentLen = utils.LenRune(*bodyStruct.Content)
		var isEmpty = contentLen < 1
		return !isEmpty, contentLen
	}

	// body has content field?
	var isContent, contentLength = checkContent()

	// pre check content.
	if contentLength > 816000 {
		validationErr.New("content")("max length is 816000")
		err = &validationErr
		return
	}

	// body has is_published field?
	var isPublished = checkIsPublished()
	// body has cover_id field?
	var isCoverID = checkCoverID()

	// check request method.
	switch requestMethod {
	default:
		validationErr.New("request method")("wrong request method")
		err = &validationErr
		return
	// POST (create) = need minimal body.
	case http.MethodPost:
		var invalid = !isContent
		if invalid {
			validationErr.New("request method")("for POST request expected at least 'content' field")
			err = &validationErr
			return
		}
	// PUT (update full) = need full body.
	case http.MethodPut:
		var invalid = !isPublished || !isCoverID || !isTitle || !isContent
		if invalid {
			validationErr.New("request method")("for PUT method expected all fields")
			err = &validationErr
			return
		}
	// PATCH (update) = need at least one field.
	case http.MethodPatch:
		var invalid = !isPublished && !isCoverID && !isTitle && !isContent
		if invalid {
			validationErr.New("request method")("for PATCH method expected at least one field")
			err = &validationErr
			return
		}
	}

	if isChangingMode {
		var articleCopy = *reference
		filtered = &articleCopy
	} else {
		filtered = &Model{}
	}

	// check body fields / add to filtered.

	if isPublished {
		filtered.IsPublished = *bodyStruct.IsPublished
	} else {
		filtered.IsPublished = false
	}

	if isCoverID {
		filtered.CoverID = bodyStruct.CoverID
	}

	if isTitle {
		filtered.Title = *bodyStruct.Title
	}

	if isContent {
		filtered.Content = *bodyStruct.Content
	}

	// check filtered fields / filter filtered (¯\_(ツ)_/¯).

	// cover.
	if filtered.CoverID != nil {
		var file = file.Model{}
		file.ID = *filtered.CoverID
		var isFound bool
		isFound, err = file.FindByID()
		// DB error.
		if err != nil {
			return
		}
		var resetCover = false
		// file does not exists.
		if !isFound {
			resetCover = true
		} else {
			// can this file be a cover?
			var extension = strings.ToUpper(file.Extension)
			var isCoverable = extension == "JPG" || extension == "JPEG" || extension == "PNG" ||
				extension == "GIF" || extension == "WEBP" || extension == "MP4"
			// validation error: bad article cover.
			if !isCoverable {
				resetCover = true
			}
		}
		if resetCover {
			filtered.CoverID = nil
		}
	}

	return
}
