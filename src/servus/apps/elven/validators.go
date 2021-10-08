package elven

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"servus/core/modules/errorMan"
	"servus/core/modules/validator"
	"strconv"
	"strings"
)


// validatorUsername - validate username from ModelUser.
func validatorUsername(username string) error {
	if len(username) < 4 || len(username) > 24 {
		return errors.New("username: min length 4 and max 24")
	}
	if !validator.IsAlphanumeric(&username) {
		return errors.New("username: allowed only alphanumeric")
	}
	return nil
}

// validatorPassword - validate ModelUser password.
func validatorPassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return errors.New("password: min length 8 and max 64")
	}
	if !validator.IsAlphanumericWithSymbols(&password) {
		return errors.New("password: allowed only alphanumeric and some symbols")
	}
	return nil
}

func (a *entityArticle) validatorControllerUpdateOne(request *http.Request) (body BodyArticle, em *errorMan.EValidation, err error){
	em = errorMan.NewValidation()
	bodyChange := &BodyArticle{}
	err = json.NewDecoder(request.Body).Decode(bodyChange)
	if err != nil {
		em.Add("body", "wrong request body provided.")
		return
	}
	if len(body.Title) > 124 {
		em.Add("title", "max length 124.")
	}
	return
}

// validatorControllerGetAll - validate query params in request depending to ModelArticle
// if validation error - returns errorMan JSON (err.Error()).
func (a *entityArticle) validatorControllerGetAll(request *http.Request, isAdmin bool) (val queryArticleControllerGetAll, em *errorMan.EValidation, err error) {
	val = queryArticleControllerGetAll{}
	em = errorMan.NewValidation()
	var queryParams = request.URL.Query()
	// validate "show" param
	var show = queryParams.Get("show")
	if len(show) == 0 {
		show = "published"
	} else {
		strings.ToLower(show)
	}
	var isShowInvalid = show != "published" && show != "drafts" && show != "all"
	var isShowForbidden = (show == "drafts" || show == "all") && !isAdmin
	if isShowInvalid || isShowForbidden {
		em.Add("show", "wrong value provided.")
	} else {
		val.show = show
	}
	// validate by param.
	var by = queryParams.Get("by")
	if len(by) == 0 {
		by = "published"
	} else {
		strings.ToLower(by)
	}
	var isByInvalid = by != "created" && by != "published" && by != "updated"
	var isByForbidden = (by == "updated" || by == "created") && !isAdmin
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
	// validate "cursor" param
	var cursor = queryParams.Get("cursor")
	if len(cursor) == 0 {
		cursor = "0"
	}
	val.cursor = cursor
	// finally
	return
}
