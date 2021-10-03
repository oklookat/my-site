package elven

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core/modules/errorCollector"
	"servus/core/modules/validator"
	"strconv"
	"strings"
)


type validatedArticlesGetAll struct {
	cursor  string
	show    string
	by      string
	start   string
	preview bool
}

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

// validatorAuth - validate request body in auth.
func validatorAuth(username string, password string, authType string) error {
	var ec = errorCollector.New()
	if validator.IsEmpty(&username) {
		ec.AddEValidationEmpty([]string{"username"})
	}
	if validator.IsEmpty(&password) {
		ec.AddEValidationEmpty([]string{"password"})
	}
	if validator.IsEmpty(&authType) {
		ec.AddEValidationEmpty([]string{"authType"})
	} else {
		var isAuthType = authType == "cookie" || authType == "direct"
		if !isAuthType {
			ec.AddEValidationAllowed([]string{"type"}, []string{"cookie", "direct"})
		}
	}
	if ec.HasErrors() {
		return errors.New(ec.GetErrors())
	}
	return nil
}

// validatorArticlesGetAll - validate query params in request depending to ModelArticle
// if validation error - returns errorCollector JSON (err.Error()).
func validatorArticlesGetAll(request *http.Request, isAdmin bool) (validated validatedArticlesGetAll, err error) {
	validated = validatedArticlesGetAll{}
	var ec = errorCollector.New()
	var queryParams = request.URL.Query()
	// validate "show" param
	var show = queryParams.Get("show")
	if show == "" {
		show = "published"
	} else {
		strings.ToLower(show)
	}
	var isShowInvalid = show != "published" && show != "drafts" && show != "all"
	switch isShowInvalid {
	case true:
		ec.AddEValidationAllowed([]string{"show"}, []string{"published", "drafts", "all"})
		break
	case false:
		if (show == "drafts" || show == "all") && !isAdmin {
			ec.AddEAuthForbidden([]string{"show"})
		}
		break
	}
	validated.show = show
	// validate "by" param
	var by = queryParams.Get("by")
	if by == "" {
		show = "published"
	} else {
		strings.ToLower(show)
	}
	var isByInvalid = by != "created" && by != "published" && by != "updated"
	if isByInvalid {
		ec.AddEValidationAllowed([]string{"by"}, []string{"created", "published", "updated"})
	} else {
		if by == "created" || by == "updated" {
			if !isAdmin {
				ec.AddEAuthForbidden([]string{"by"})
			}
		}
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
	validated.by = by
	// validate "start" param
	var start = queryParams.Get("start")
	if start == "" {
		start = "newest"
	} else {
		strings.ToLower(start)
	}
	var isStartInvalid = start != "newest" && start != "oldest"
	if isStartInvalid {
		ec.AddEValidationAllowed([]string{"start"}, []string{"newest", "oldest"})
	} else {
		switch start {
		case "newest":
			start = "DESC"
			break
		case "oldest":
			start = "ASC"
			break
		}
	}
	validated.start = start
	// validate "preview" param
	var preview = queryParams.Get("preview")
	if preview == "" {
		preview = "true"
	}
	var previewBool bool
	previewBool, err = strconv.ParseBool(preview)
	if err != nil {
		ec.AddEValidationAllowed([]string{"preview"}, []string{"boolean"})
		previewBool = true
	}
	validated.preview = previewBool
	// validate "cursor" param
	var cursor = queryParams.Get("cursor")
	if cursor == "" {
		cursor = "0"
	}
	validated.cursor = cursor
	// finally
	if ec.HasErrors() {
		return validated, errors.New(ec.GetErrors())
	}
	return validated, nil
}

func validatorArticlesPost(body *ControllerArticlesPostBody) error{
	var ec = errorCollector.New()
	if len(body.Title) > 124 {
		ec.AddEValidationMinMax([]string{"title"}, 1, 124)
	}
	if len(body.Title) == 0 {
		body.Title = "Без названия"
	}
	//if len(body.Content) > 512000 {
	//	ec.AddEValidationMinMax([]string{"title"}, 1, 512000)
	//}
	if ec.HasErrors() {
		return errors.New(ec.GetErrors())
	}
	return nil
}
