package base

type RequestError interface {
	// Server - 500 error.
	Server() string
	// Get500 - 401 error.
	NotAuthorized() string
	// Get500 - 403 error.
	Forbidden() string
	// Get500 - 404 error.
	NotFound() string
}

// Validate - create validator.
type Validate interface {
	Create() Validator
}

// Validator - provides functions to validate things.
type Validator interface {
	Add(field string)
	HasErrors() bool
	GetJSON() string
}
