package errorCollector

import (
	"encoding/json"
)

type ErrorsSlice []interface{}
type ErrorCollector struct {
	Errors struct {
		ErrorsSlice `json:"errors"`
	}
}

// HasErrors - check is errorCollector has errors.
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.Errors.ErrorsSlice) > 0
}

// GetErrors - get errors in JSON.
func (ec *ErrorCollector) GetErrors() string {
	var bytes, _ = json.Marshal(ec.Errors)
	return string(bytes)
}

// Clear - remove all errors.
func (ec *ErrorCollector) Clear() {
	ec.Errors.ErrorsSlice = nil
}

// AddEUnknown - like 500 error.
func (ec *ErrorCollector) AddEUnknown(issuers []string, message string) {
	err := EUnknown{
		EError{StatusCode: 500, ErrorCode: "E_UNKNOWN", Issuers: issuers, Message: message},
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddECustom - like 'I need show message for users'.
func (ec *ErrorCollector) AddECustom(issuers []string, message string, data interface{}) {
	err := ECustom{
		EError{StatusCode: 500, ErrorCode: "E_CUSTOM", Issuers: issuers, Message: message},
		data,
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddEAuthIncorrect - like 'wrong username or password'.
func (ec *ErrorCollector) AddEAuthIncorrect(issuers []string) {
	err := EAuthIncorrect{
		EError{StatusCode: 403, ErrorCode: "E_AUTH_INCORRECT", Issuers: issuers, Message: "Wrong credentials."},
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddEAuthForbidden - like 'only admins'.
func (ec *ErrorCollector) AddEAuthForbidden(issuers []string) {
	err := EAuthForbidden{
		EError{StatusCode: 403, ErrorCode: "E_AUTH_FORBIDDEN", Issuers: issuers, Message: "Access denied."},
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddENotFound - like 'article not found'.
func (ec *ErrorCollector) AddENotFound(issuers []string) {
	err := ENotFound{
		EError{StatusCode: 404, ErrorCode: "E_NOTFOUND", Issuers: issuers, Message: "Not found."},
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddEValidationMinMax - like 'min length for username is 4 symbols'.
func (ec *ErrorCollector) AddEValidationMinMax(issuers []string, min int, max int) {
	err := EValidationMinMax{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_MINMAX", Issuers: issuers, Message: "Too many or not enough characters."},
		min, max,
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddEValidationEmpty - like 'title cannot be empty'.
func (ec *ErrorCollector) AddEValidationEmpty(issuers []string) {
	err := EValidationEmpty{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_EMPTY", Issuers: issuers, Message: "These things cannot be empty."},
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddEValidationAllowed - like 'allowed only numbers' or 'is_published must be true or false'.
func (ec *ErrorCollector) AddEValidationAllowed(issuers []string, allowed []string) {
	err := EValidationAllowed{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_ALLOWED", Issuers: issuers, Message: "You send some kind of nonsense."},
		allowed,
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}

// AddEValidationInvalid - like 'request contains file, but file broken'.
func (ec *ErrorCollector) AddEValidationInvalid(issuers []string, message string) {
	err := EValidationInvalid{
		EError{StatusCode: 400, ErrorCode: "E_VALIDATION_INVALID", Issuers: issuers, Message: message},
	}
	ec.Errors.ErrorsSlice = append(ec.Errors.ErrorsSlice, err)
}
