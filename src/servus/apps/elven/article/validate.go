package article

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"servus/apps/elven/base"
	"servus/apps/elven/model"
	"strconv"
	"strings"
)

// validate params to get paginated articles.
func ValidateGetParams(a *base.ArticleGetParams, params url.Values, isAdmin bool) (val base.Validator) {
	val = validate.Create()
	// "show" param
	var show = params.Get("show")
	if len(show) == 0 {
		show = "published"
	}
	var isShowPublished = show == "published"
	var isShowDrafts = show == "drafts"
	var isShowInvalid = !isShowPublished && !isShowDrafts
	var isShowForbidden = isShowDrafts && !isAdmin
	if isShowInvalid || isShowForbidden {
		val.Add("show")
	} else {
		a.Show = show
	}
	// "by" param.
	var by = params.Get("by")
	if len(by) == 0 {
		by = "published"
	}
	switch by {
	default:
		val.Add("by")
	case "created":
		if !isAdmin {
			val.Add("by")
		}
		by = "created_at"
	case "updated":
		if !isAdmin {
			val.Add("by")
		}
		by = "updated_at"
	case "published":
		by = "published_at"
	}
	a.By = by
	// "start" param.
	var start = params.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		val.Add("start")
	} else {
		if isNewest {
			start = "DESC"
		} else if isOldest {
			start = "ASC"
		}
		a.Start = start
	}
	// "preview" param.
	var preview = params.Get("preview")
	if len(preview) == 0 {
		preview = "true"
	}
	var previewBool bool
	previewBool, err := strconv.ParseBool(preview)
	if err != nil {
		val.Add("preview")
	}
	a.Preview = previewBool
	// "page" param
	var pageStr = params.Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		val.Add("page")
	} else {
		a.Page = page
	}
	// "without category" param
	var withoutCategory = params.Get("without_category")
	a.WithoutCategory = len(withoutCategory) > 0 && withoutCategory == "true"
	if !a.WithoutCategory {
		// "category name" param
		var categoryName = params.Get("category_name")
		if len(categoryName) > 0 {
			a.CategoryName = &categoryName
		}
	}
	return
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
func ValidateBody(requestMethod string, body io.ReadCloser, reference *model.Article) *model.Article {

	// updating existing article?
	var isChangingMode = requestMethod == http.MethodPut || requestMethod == http.MethodPatch
	// no validate if no reference in changing mode.
	if isChangingMode && reference == nil {
		return nil
	}

	// decode body.
	var bodyStruct = &base.ArticleBody{}
	err := json.NewDecoder(body).Decode(bodyStruct)
	if err != nil {
		return nil
	}

	// FUNC title.
	//
	// returns: is valid, title length.
	var checkTitle = func() (bool, int) {
		var isNil = bodyStruct.Title != nil
		if isNil {
			return false, 0
		}
		var titleLength = call.Utils.LenRune(*bodyStruct.Title)
		var isEmpty = titleLength < 1
		return !isEmpty, titleLength
	}

	// body has title field?
	var isTitle, titleLength = checkTitle()

	// pre check title.
	if isTitle && titleLength > 124 {
		return nil
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
		var contentLen = call.Utils.LenRune(*bodyStruct.Content)
		var isEmpty = contentLen < 1
		return !isEmpty, contentLen
	}

	// body has content field?
	var isContent, contentLength = checkContent()

	// pre check content.
	if contentLength > 256000 {
		return nil
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
		return nil
	// POST (create) = need minimal body.
	case http.MethodPost:
		var invalid = !isContent
		if invalid {
			return nil
		}
	// PUT (update full) = need full body.
	case http.MethodPut:
		var invalid = !isPublished || !isCategoryID || !isCoverID || !isTitle || !isContent
		if invalid {
			return nil
		}
	// PATCH (update) = need at least one field.
	case http.MethodPatch:
		var invalid = !isPublished && !isCategoryID && !isCoverID && !isTitle && !isContent
		if invalid {
			return nil
		}
	}

	var filtered *model.Article
	if isChangingMode {
		var articleCopy = *reference
		filtered = &articleCopy
	} else {
		filtered = &model.Article{}
	}

	// check body fields / add to filtered.

	if isPublished {
		filtered.IsPublished = *bodyStruct.IsPublished
	}
	if isCategoryID {
		if *bodyStruct.CategoryID == "nope" {
			bodyStruct.CategoryID = nil
		}
		filtered.CategoryID = bodyStruct.CategoryID
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
			return nil
		}
	}
	// cover.
	if filtered.CoverID != nil {
		var file = model.File{}
		file.ID = *filtered.CoverID
		found, err := file.FindByID()
		// validation error: fake file or DB error.
		if err != nil || !found {
			return nil
		}
		// can this file be a cover?
		// TODO: write check.
	}

	return filtered
}

func ValidateCategoryBody(c *base.CategoryBody, body io.ReadCloser) (val base.Validator) {
	val = validate.Create()
	// body.
	err := json.NewDecoder(body).Decode(c)
	if err != nil {
		val.Add("body")
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
	var nameLen = call.Utils.LenRune(c.Name)
	var notValid = nameLen < 1 || nameLen > 24
	if notValid {
		val.Add("name")
	}
	return
}
