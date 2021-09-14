package errorCollector

type EError struct {
	StatusCode int      `json:"statusCode"`
	ErrorCode  string   `json:"errorCode"`
	Issuers    []string `json:"issuers"`
	Message    string   `json:"message"`
}

// EUnknown like 500 error
type EUnknown struct {
	EError
}

// ECustom like 'I need show message for users'
type ECustom struct {
	EError
	Data interface{} `json:"data"`
}

// EAuthIncorrect like 'wrong username or password'
type EAuthIncorrect struct {
	EError
}

// EAuthForbidden like 'only admins'
type EAuthForbidden struct {
	EError
}

// ENotFound like 'article not found'
type ENotFound struct {
	EError
}

// EValidationMinMax like 'min length for username is 4 symbols'
type EValidationMinMax struct {
	EError
	Min int `json:"min"`
	Max int `json:"max"`
}

// EValidationEmpty like 'title cannot be empty'
type EValidationEmpty struct {
	EError
}

// EValidationAllowed like 'allowed only numbers' or 'is_published must be true or false'
type EValidationAllowed struct {
	EError
	Allowed []string `json:"allowed"`
}

// EValidationInvalid like 'request contains file, but file broken'
type EValidationInvalid struct {
	EError
}

