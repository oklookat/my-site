package errorCollector

import (
	"encoding/json"
	"fmt"
)


func New() *errorCollector {
	return &errorCollector{}
}

func (ec *errorCollector) HasErrors() bool {
	return len(ec.errorsArray) > 0
}

func (ec *errorCollector) GetErrors() string {
	var bytes, _ = json.Marshal(ec.errorsArray)
	return fmt.Sprintf(`{"errors": %v}`, string(bytes))
}

// AddEUnknown like 500 error
func (ec *errorCollector) AddEUnknown(issuers []string, message string) {
	_error := EUnknown{
		EError{StatusCode: 500, ErrorCode: "E_UNKNOWN", Issuers: issuers, Message: message},
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddECustom like 'I need show message for users'
func (ec *errorCollector) AddECustom(issuers []string, message string, data interface{}) {
	_error := ECustom{
		EError{StatusCode: 500, ErrorCode: "E_CUSTOM", Issuers: issuers, Message: message},
		data,
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddEAuthIncorrect like 'wrong username or password'
func (ec *errorCollector) AddEAuthIncorrect(issuers []string) {
	_error := EAuthIncorrect{
		EError{StatusCode: 403, ErrorCode: "E_AUTH_INCORRECT", Issuers: issuers, Message: "Wrong credentials."},
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddEAuthForbidden like 'only admins'
func (ec *errorCollector) AddEAuthForbidden(issuers []string) {
	_error := EAuthForbidden{
		EError{StatusCode: 403, ErrorCode: "E_AUTH_FORBIDDEN", Issuers: issuers, Message: "Access denied."},
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddENotFound like 'article not found'
func (ec *errorCollector) AddENotFound(issuers []string) {
	_error := ENotFound{
		EError{StatusCode: 404, ErrorCode: "E_NOTFOUND", Issuers: issuers, Message: "Not found."},
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddEValidationMinMax like 'min length for username is 4 symbols'
func (ec *errorCollector) AddEValidationMinMax(issuers []string, min int, max int) {
	_error := EValidationMinMax{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_MINMAX", Issuers: issuers, Message: "Too many or not enough characters."},
		min, max,
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddEValidationEmpty like 'title cannot be empty'
func (ec *errorCollector) AddEValidationEmpty(issuers []string) {
	_error := EValidationEmpty{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_EMPTY", Issuers: issuers, Message: "These things cannot be empty."},
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddEValidationAllowed like 'allowed only numbers' or 'is_published must be true or false'
func (ec *errorCollector) AddEValidationAllowed(issuers []string, allowed []string) {
	_error := EValidationAllowed{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_ALLOWED", Issuers: issuers, Message: "These things not allowed."},
		allowed,
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}

// AddEValidationInvalid like 'request contains file, but file broken'
func (ec *errorCollector) AddEValidationInvalid(issuers []string, message string) {
	_error := EValidationInvalid{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_INVALID", Issuers: issuers, Message: message},
	}
	ec.errorsArray = append(ec.errorsArray, _error)
}
