package elven

// BootApp - start app.
func BootApp() {
	bootEntities()
	oCmd.boot()
	bootRoutes()
}
