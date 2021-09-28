package elven

import (
	"net/http"
	"servus/core"
	"servus/core/modules/routerica"
)

func bootRoutes(){
	router := routerica.New()
	router.Use(core.Middleware.MiddlewareCORS, core.Middleware.MiddlewareSecurity, core.Middleware.MiddlewareAsJSON)
	var elvenRoutes = router.Group("/elven")
	elvenRoutes.POST("/auth/login", controllerAuthLogin)
	elvenRoutes.POST("/auth/logout", controllerAuthLogout)
	http.Handle("/", router)
}
