package user

import (
	"encoding/json"
	"io"
	"servus/apps/elven/base"
)

type Body struct {
	What     string
	Password string
	NewValue string
}

func (u *Body) Validate(body io.ReadCloser) base.Validator {
	val := validate.Create()
	err := json.NewDecoder(body).Decode(u)
	if err != nil {
		val.Add("body")
	}
	return val
}
