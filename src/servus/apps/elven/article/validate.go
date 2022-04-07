package article

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
	"servus/core/external/utils"
	"strconv"
	"strings"
)

// validate params to get paginated articles.
func ValidateGetParams(a *base.ArticleGetParams, params url.Values, isAdmin bool) {
	// "show" param
	var show = params.Get("show")
	var isShowPublished = strings.EqualFold(show, "published")
	var isShowDrafts = strings.EqualFold(show, "drafts")
	var isShowInvalid = !isShowPublished && !isShowDrafts
	var isShowForbidden = isShowDrafts && !isAdmin
	if isShowPublished || isShowInvalid || isShowForbidden {
		show = "published"
	}
	a.Show = show

	// "by" param.
	var by = params.Get("by")
	var isByUpdated = strings.EqualFold(by, "updated")
	var isByCreated = strings.EqualFold(by, "created")
	var isByPublished = strings.EqualFold(by, "published")
	var isByInvalid = !isByUpdated && !isByCreated && !isByPublished
	var isByForbidden = !isAdmin && (isByCreated || isByUpdated)
	if isByPublished || isByInvalid || isByForbidden {
		by = "published"
	}
	by = by + "_at"
	a.By = by

	// "start" param.
	var start = params.Get("start")
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isNewest || isStartInvalid {
		start = "DESC"
	} else if isOldest {
		start = "ASC"
	}
	a.Start = start

	// "preview" param.
	preview, err := strconv.ParseBool(params.Get("preview"))
	if err != nil {
		preview = true
	}
	a.Preview = preview

	// "page" param
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	a.Page = page

	// "without category" param
	a.WithoutCategory, err = strconv.ParseBool(params.Get("without_category"))
	if err != nil {
		a.WithoutCategory = true
	}
	if !a.WithoutCategory {
		// "category name" param
		var categoryName = params.Get("category_name")
		if len(categoryName) > 0 {
			a.CategoryName = &categoryName
		} else {
			a.WithoutCategory = true
		}
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
func ValidateBody(requestMethod string, body io.ReadCloser, reference *model.Article) (filtered *model.Article, err error) {

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
	var bodyStruct = &base.ArticleBody{}
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

	// FUNC category id.
	var checkCategoryID = func() bool {
		return bodyStruct.CategoryID != nil
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
	if contentLength > 256000 {
		validationErr.New("content")("max length is 256000")
		err = &validationErr
		return
	}

	// body has is_published field?
	var isPublished = checkIsPublished()
	// body has category_id field?
	var isCategoryID = checkCategoryID()
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
		var invalid = !isPublished || !isCategoryID || !isCoverID || !isTitle || !isContent
		if invalid {
			validationErr.New("request method")("for PUT method expected all fields")
			err = &validationErr
			return
		}
	// PATCH (update) = need at least one field.
	case http.MethodPatch:
		var invalid = !isPublished && !isCategoryID && !isCoverID && !isTitle && !isContent
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
		filtered = &model.Article{}
	}

	// check body fields / add to filtered.

	if isPublished {
		filtered.IsPublished = *bodyStruct.IsPublished
	} else {
		filtered.IsPublished = false
	}

	if isCategoryID {
		if *bodyStruct.CategoryID == "nope" {
			bodyStruct.CategoryID = nil
		}
		filtered.CategoryID = bodyStruct.CategoryID
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

	// category.
	if filtered.CategoryID != nil {
		var cat = model.ArticleCategory{}
		cat.ID = *filtered.CategoryID
		found, err := cat.FindByID()
		// validation error: fake category or DB error.
		if err != nil || !found {
			// reset category.
			filtered.CategoryID = nil
		}
	}

	// cover.
	if filtered.CoverID != nil {
		var file = model.File{}
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

func ValidateCategoryBody(c *base.CategoryBody, body io.ReadCloser) (err error) {
	var validationErr = base.ValidationError{}

	// body.
	if err = json.NewDecoder(body).Decode(c); err != nil {
		return
	}

	// replace new lines with one space
	reg, _ := regexp.Compile(`[\r\n]`)
	c.Name = reg.ReplaceAllString(c.Name, " ")

	// replace 2+ spaces with one space
	reg, _ = regexp.Compile(`[^\S]{2,}`)
	c.Name = reg.ReplaceAllString(c.Name, "")

	// remove spaces at start and end
	c.Name = strings.Trim(c.Name, " ")

	// length?
	var nameLen = utils.LenRune(c.Name)
	var notValid = nameLen < 1 || nameLen > 24
	if notValid {
		validationErr.New("name")("min length is 1 and max is 24")
		err = &validationErr
		return
	}

	return
}
