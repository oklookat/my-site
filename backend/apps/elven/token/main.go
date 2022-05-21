package token

import (
	"errors"
	"servus/core"
)

var call *core.Instance

type Starter struct {
	Core *core.Instance
}

func Start(s *Starter) error {
	// check.
	if s == nil {
		return errors.New("starter nil pointer")
	}
	if s.Core == nil {
		return errors.New("core nil pointer")
	}

	// set.
	call = s.Core

	return nil
}
