package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
)

func bootRoutes(){
	// set up router
	router := mux.NewRouter()
	router.Use(core.Middleware.MiddlewareSecurity, core.Middleware.MiddlewareAsJSON)
	routerElven := router.PathPrefix("/elven/").Subrouter()
	//routerAdminOnly := routerElven.PathPrefix("").Subrouter()
	//routerAdminOnly.Use(middlewareAdminOnly)
	// handlers
	routerElven.HandleFunc("/auth/login", controllerAuthLogin).Methods("POST")
	routerElven.HandleFunc("/auth/logout", controllerAuthLogout).Methods("POST")
	//routerSub.HandleFunc("/articles", elControllers.Articles)
	//routerSub.HandleFunc("/articles", elControllers.Files)

	// use very important middleware before router
	var useCORS = core.Middleware.MiddlewareCORS(router)
	http.Handle("/", useCORS)
}
