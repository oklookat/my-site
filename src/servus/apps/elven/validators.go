package elven

import (
	"encoding/json"
	"net/http"
	"servus/core/modules/errorMan"
	"servus/core/modules/validator"
	"strconv"
	"strings"
)


// validatorControllerLogin - validate request body when user try to log in.
func (a *entityAuth) validatorControllerLogin(request *http.Request) (val *bodyAuth, em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	val = &bodyAuth{}
	err = json.NewDecoder(request.Body).Decode(&val)
	if err != nil {
		em.Add("body", "wrong body provided.")
		return
	}
	var username = val.Username
	var password = val.Password
	var authType = val.Type
	if validator.IsEmpty(&username) {
		em.Add("username", "cannot be empty.")
	}
	if validator.IsEmpty(&password) {
		em.Add("password", "cannot be empty.")
	}
	if validator.IsEmpty(&authType) {
		em.Add("type", "cannot be empty.")
	} else {
		var isAuthType = authType == "cookie" || authType == "direct"
		if !isAuthType {
			em.Add("type", "wrong type.")
		}
	}
	return
}

// validatorControllerGetAll - validate query params when getting articles list.
func (a *entityArticle) validatorControllerGetAll(request *http.Request, isAdmin bool) (val queryArticleGetAll, em *errorMan.EValidation, err error) {
	val = queryArticleGetAll{}
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
		break
	case "updated":
		by = "updated_at"
		break
	case "published":
		by = "published_at"
		break
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

// validatorBody - validate request body when POST or PUT.
func (a *entityArticle) validatorBody(request *http.Request) (val *BodyArticle, em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	val = &BodyArticle{}
	err = json.NewDecoder(request.Body).Decode(val)
	if err != nil {
		em.Add("title", "wrong value provided.")
		em.Add("content", "wrong value provided.")
		return
	}
	if len(val.Title) > 124 {
		em.Add("title", "max length 124.")
	}
	return
}

// validatorControllerGetAll - validate query params when getting files list.
func (f *entityFile) validatorControllerGetAll(request *http.Request, isAdmin bool) (val queryFileGetAll, em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	val = queryFileGetAll{}
	var queryParams = request.URL.Query()
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
	}
	val.start = start
	// validate cursor param.
	var cursor = queryParams.Get("cursor")
	if len(cursor) == 0 {
		cursor = "0"
	}
	val.cursor = cursor
	// validate by param.
	var by = queryParams.Get("by")
	if len(by) == 0 {
		by = "created"
	}
	var isByCreated = strings.EqualFold(by, "created")
	var isByInvalid = !isByCreated
	var isByForbidden = (isByCreated) && !isAdmin
	if isByInvalid || isByForbidden {
		em.Add("by", "wrong value provided.")
	}
	switch by {
	case "created":
		by = "created_at"
		break
	}
	val.by = by
	return
}
