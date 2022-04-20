package base

import (
	"fmt"
)

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
