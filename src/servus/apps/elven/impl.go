package elven

import (
	"servus/apps/elven/base"
	"servus/core/external/errorMan"
)

type Validate struct {
}

func (v *Validate) Create() base.Validator {
	var val = errorMan.Validate{}
	return val.Create()
}
