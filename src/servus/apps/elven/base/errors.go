package base

import (
	"fmt"
)

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

type ValidationError struct {
	Issuer  string
	Message string
}

func (v *ValidationError) New(issuer string) func(message string) {
	v.Issuer = issuer
	return func(message string) {
		v.Message = message
	}
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("[validation/%v] %v", v.Issuer, v.Message)
}
