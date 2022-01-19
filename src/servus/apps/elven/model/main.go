package model

import (
	"servus/core"
	"servus/core/external/database"
)

var call *core.Instance
var IntAdapter = database.Adapter[int]{}
var StringAdapter = database.Adapter[string]{}

func Boot(c *core.Instance) {
	call = c
}
