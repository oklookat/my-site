package base

type RequestError interface {
	// 500 error.
	Server() string
	// 401 error.
	NotAuthorized() string
	// 403 error.
	Forbidden() string
	// 404 error.
	NotFound() string
	// 409 error.
	Exists() string
}

// create validator.
type Validate interface {
	Create() Validator
}

// provides functions to validate things.
type Validator interface {
	Add(field string)
	HasErrors() bool
	GetJSON() string
}
