package errorMan

import (
	"encoding/json"
)

// NewValidation - new validation error collector (400).
func NewValidation() *EValidation {
	err := &EValidation{
		eError{StatusCode: 400, ErrorCode: "E_BAD_REQUEST"},
		make(map[string]string, 0),
	}
	return err
}

// Add - add field into EValidation.InvalidFields.
func (e *EValidation) Add(field string, message string) {
	e.InvalidFields[field] = message
}

// HasErrors - check is validation errors exists.
func (e *EValidation) HasErrors() bool {
	return len(e.InvalidFields) > 0
}

// GetJSON - get validation errors in JSON format.
func (e *EValidation) GetJSON() string {
	var bytes, _ = json.Marshal(e)
	return string(bytes)
}
