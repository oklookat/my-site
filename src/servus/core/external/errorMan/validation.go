package errorMan

import (
	"encoding/json"
)

type Validate struct {
}

func (v Validate) Create() *Validator {
	var val = &Validator{}
	val.eError = eError{
		StatusCode: 400,
		ErrorCode:  "E_BAD_REQUEST",
	}
	val.InvalidFields = make([]string, 0)
	return val
}

type Validator struct {
	eError
	InvalidFields []string `json:"invalid_fields"`
}

// add field into EValidation.InvalidFields.
func (e *Validator) Add(field string) {
	e.InvalidFields = append(e.InvalidFields, field)
}

// check is validation errors exists.
func (e Validator) HasErrors() bool {
	return len(e.InvalidFields) > 0
}

// get validation errors in JSON format.
func (e Validator) GetJSON() string {
	var bytes, _ = json.Marshal(e)
	return string(bytes)
}
