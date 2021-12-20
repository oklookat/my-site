package pipe

import "servus/core"

var call *core.Instance

func Boot(c *core.Instance) {
	call = c
}
