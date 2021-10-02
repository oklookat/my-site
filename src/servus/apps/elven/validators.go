package elven

import (
	"github.com/pkg/errors"
	"net/http"
	"servus/core/modules/errorCollector"
	"servus/core/modules/validator"
	"strconv"
	"strings"
)

//func validatorUserReg(username string, password string) (bool, string) {
//	var ec = errorCollector.New()
//	// username and password start //
//	var usernameIssuer = []string{"username"}
//	var passwordIssuer = []string{"password"}
//	var failed = 0
//	if len(username) < 4 || len(username) > 24 {
//		ec.AddEValidationMinMax(usernameIssuer, 4, 24)
//		failed++
//	}
//	if len(password) < 8 || len(password) > 64 {
//		ec.AddEValidationMinMax(passwordIssuer, 8, 64)
//		failed++
//	}
//	if failed == 2 { // if username and password failed
//		return true, ec.GetErrors()
//	}
//	isAlphaNumeric := regexp.MustCompile("^[a-zA-Z0-9_]*$")
//	if !isAlphaNumeric.MatchString(username) {
//		ec.AddEValidationAllowed(usernameIssuer, []string{"alphanumeric"})
//	}
//	isPassword := regexp.MustCompile("^[A-Za-z0-9_@!./#&+*%-]*$")
//	if !isPassword.MatchString(password) {
//		ec.AddEValidationAllowed(passwordIssuer, []string{"alphanumeric-and-some-symbols"})
//	}
//	// username and password end //
//	var role = user.Role
//	var roleIssuer = []string{"role"}
//	if len(role) > 0 {
//		if role != "user" && role != "admin" {
//			ec.AddEValidationAllowed(roleIssuer, []string{"user", "admin"})
//		}
//	}
//	var regIP = user.RegIP
//	if len(regIP) > 0 {
//		err := pValidate.Var(regIP, "ip")
//		if err != nil {
//			ec.AddEValidationInvalid([]string{"RegIP"}, "invalid IP address")
//		}
//	}
//
//	//var regAgent = user.RegAgent
//	if ec.HasErrors() {
//		println(ec.GetErrors())
//		return true, ec.GetErrors()
//	}
//	return false, ""
//}

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

type validatedArticleQuery struct {
	cursor  string
	show    string
	by      string
	start   string
	preview bool
}

// validatorArticleQueryParams - validate query params in request depending to ModelArticle
// if validation error - returns errorCollector JSON (err.Error()).
func validatorArticleQueryParams(request *http.Request, isAdmin bool) (validatedParams validatedArticleQuery, err error) {
	isAdmin = true
	validatedParams = validatedArticleQuery{}
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
	validatedParams.show = show
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
	validatedParams.by = by
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
	validatedParams.start = start
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
	validatedParams.preview = previewBool
	// validate "cursor" param
	var cursor = queryParams.Get("cursor")
	if cursor == "" {
		cursor = "0"
	}
	validatedParams.cursor = cursor
	// finally
	if ec.HasErrors() {
		return validatedParams, errors.New(ec.GetErrors())
	}
	return validatedParams, nil
}
