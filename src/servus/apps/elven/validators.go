package elven

import (
	"github.com/pkg/errors"
	"regexp"
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
	isAlphaNumeric := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !isAlphaNumeric.MatchString(username) {
		return errors.New("username: allowed only alphanumeric")
	}
	return nil
}

func validatorPassword(password string) error {
	if len(password) < 8 || len(password) > 64 {
		return errors.New("password: min length 8 and max 64")
	}
	isPassword := regexp.MustCompile("^[A-Za-z0-9_@!./#&+*%-]*$")
	if !isPassword.MatchString(password) {
		return errors.New("password: allowed only alphanumeric and some symbols")
	}
	return nil
}
