package elven

import "servus/core"

var instance *core.Core

// BootApp - start app.
func BootApp(c *core.Core) {
	instance = c
	bootEntities()
	oCmd.boot()
	bootRoutes()
}
