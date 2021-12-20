package article

import (
	"encoding/json"
	"net/http"
	"servus/core/external/errorMan"
	"strconv"
	"strings"
)

// validatorControllerGetAll - validate query params when getting articles list.
func (a *validator) controllerGetAll(request *http.Request, isAdmin bool) (val queryGetAll, em *errorMan.EValidation, err error) {
	val = queryGetAll{}
	em = errorMan.NewValidation()
	var queryParams = request.URL.Query()
	// validate "show" param
	var show = queryParams.Get("show")
	if len(show) == 0 {
		show = "published"
	}
	var isShowPublished = strings.EqualFold(show, "published")
	var isShowDrafts = strings.EqualFold(show, "drafts")
	var isShowInvalid = !isShowPublished && !isShowDrafts
	var isShowForbidden = isShowDrafts && !isAdmin
	if isShowInvalid || isShowForbidden {
		em.Add("show", "wrong value provided.")
	} else {
		val.show = show
	}
	// validate by param.
	var by = queryParams.Get("by")
	if len(by) == 0 {
		by = "published"
	}
	var isByCreated = strings.EqualFold(by, "created")
	var isByPublished = strings.EqualFold(by, "published")
	var isByUpdated = strings.EqualFold(by, "updated")
	var isByInvalid = !isByCreated && !isByPublished && !isByUpdated
	var isByForbidden = (isByCreated || isByUpdated) && !isAdmin
	if isByInvalid || isByForbidden {
		em.Add("by", "wrong value provided.")
	}
	switch by {
	case "created":
		by = "created_at"
	case "updated":
		by = "updated_at"
	case "published":
		by = "published_at"
	}
	val.by = by
	// validate start param.
	var start = queryParams.Get("start")
	if len(start) == 0 {
		start = "newest"
	}
	var isNewest = strings.EqualFold(start, "newest")
	var isOldest = strings.EqualFold(start, "oldest")
	var isStartInvalid = !isNewest && !isOldest
	if isStartInvalid {
		em.Add("start", "wrong value provided.")
	} else {
		if isNewest {
			start = "DESC"
		} else if isOldest {
			start = "ASC"
		}
		val.start = start
	}
	// validate preview param.
	var preview = queryParams.Get("preview")
	if len(preview) == 0 {
		preview = "true"
	}
	var previewBool bool
	previewBool, err = strconv.ParseBool(preview)
	if err != nil {
		em.Add("preview", "wrong value provided.")
		previewBool = true
	} else {
		val.preview = previewBool
	}
	// validate "page" param
	var pageStr = queryParams.Get("page")
	if len(pageStr) == 0 {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		em.Add("page", "wrong value provided.")
	} else {
		val.page = page
	}
	return
}

// TODO: when provided only one value, it make other values to null, and pass it do db.
// validatorBody - validate request body when POST or PUT.
func (a *validator) body(request *http.Request) (val *Body, em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	val = &Body{}
	err = json.NewDecoder(request.Body).Decode(val)
	if err != nil {
		em.Add("body", "wrong value provided.")
		return
	}
	var isContent = val.Content != nil
	if isContent {
		var contentInvalid = len(val.Content.Blocks) < 1
		if contentInvalid {
			em.Add("content", "wrong value provided.")
		}
	}
	var isTitle = val.Title != nil
	if isTitle {
		if len(*val.Title) < 1 {
			*val.Title = "Untitled"
		} else if len(*val.Title) > 124 {
			em.Add("title", "max length 124.")
		}
	}
	var isPublished = val.IsPublished != nil
	if !isTitle && !isContent && !isPublished {
		em.Add("body", "cannot be empty.")
	}
	return
}
