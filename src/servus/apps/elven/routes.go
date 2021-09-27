package elven

import (
	"net/http"
	"servus/core"
	"servus/core/modules/routerica"
)

func bootRoutes(){
	// set up router
	router := routerica.New()
	router.Use(core.Middleware.MiddlewareCORS, core.Middleware.MiddlewareSecurity, core.Middleware.MiddlewareAsJSON)
	var elvenRoutes = router.Group("/elven")
	elvenRoutes.POST("/auth/login", controllerAuthLogin)
	elvenRoutes.POST("/auth/logout", controllerAuthLogout)
	//routerAdminOnly := routerElven.PathPrefix("").Subrouter()
	//routerAdminOnly.Use(middlewareAdminOnly)
	// handlers
	//routerSub.HandleFunc("/articles", elControllers.Articles)
	//routerSub.HandleFunc("/articles", elControllers.Files)
	http.Handle("/", router)
}
