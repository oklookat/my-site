package elven

import (
	"github.com/pkg/errors"
	"servus/core/modules/errorCollector"
	"servus/core/modules/validator"
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

func validatorUsername(username string) error {
	if len(username) < 4 || len(username) > 24 {
		return errors.New("username: min length 4 and max 24")
	}
	if !validator.IsAlphanumeric(username) {
		return errors.New("username: allowed only alphanumeric")
	}
	return nil
}

func validatorPassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return errors.New("password: min length 8 and max 64")
	}
	if !validator.IsAlphanumericWithSymbols(password) {
		return errors.New("password: allowed only alphanumeric and some symbols")
	}
	return nil
}

func validatorAuth(username string, password string, authType string) error{
	var ec = errorCollector.New()
	if validator.IsEmpty(username) {
		ec.AddEValidationEmpty([]string{"username"})
	}
	if validator.IsEmpty(password){
		ec.AddEValidationEmpty([]string{"password"})
	}
	if validator.IsEmpty(authType){
		ec.AddEValidationEmpty([]string{"authType"})
	} else {
		var isAuthType = authType == "cookie" || authType == "direct"
		if !isAuthType{
			ec.AddEValidationAllowed([]string{"type"}, []string{"cookie", "direct"})
		}
	}
	if ec.HasErrors() {
		return errors.New(ec.GetErrors())
	}
	return nil
}
