package elven

import (
	"encoding/json"
	"net/http"
	"servus/core/external/errorMan"
	"servus/core/external/validator"
)

// validatorControllerLogin - validate request body when user try to log in.
func (a *entityAuth) validatorControllerLogin(request *http.Request) (val *bodyAuth, em *errorMan.EValidation, err error) {
	em = errorMan.NewValidation()
	val = &bodyAuth{}
	err = json.NewDecoder(request.Body).Decode(&val)
	if err != nil {
		em.Add("body", "wrong value provided.")
		return
	}
	var username = val.Username
	var password = val.Password
	var authType = val.Type
	if validator.IsEmpty(&username) {
		em.Add("username", "wrong value provided.")
	}
	if validator.IsEmpty(&password) {
		em.Add("password", "wrong value provided.")
	}
	if validator.IsEmpty(&authType) {
		em.Add("type", "wrong value provided.")
	} else {
		var isAuthType = authType == "cookie" || authType == "direct"
		if !isAuthType {
			em.Add("type", "wrong value provided.")
		}
	}
	return
}
