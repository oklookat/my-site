package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/apps/elven/elUser"
	"servus/core"
)

func bootRoutes(){
	// TODO: fix MultipartForm and add many variants to get username and password
	// TODO: test CORS and test cookie max age
	// set up router
	router := mux.NewRouter()
	router.Use(core.Middleware.MiddlewareSecurity)
	router.Use(core.Middleware.MiddlewareAsJSON)
	// handlers
	routerSub := router.PathPrefix("/elven/").Subrouter()
	routerSub.HandleFunc("/auth/login", elUser.ControllerAuthLogin).Methods("POST")
	routerSub.HandleFunc("/auth/logout", elUser.ControllerAuthLogout).Methods("POST")
	//routerSub.HandleFunc("/articles", elControllers.Articles)
	//routerSub.HandleFunc("/articles", elControllers.Files)

	// use very important middleware before router
	var useCORS = core.Middleware.MiddlewareCORS(router)
	http.Handle("/", useCORS)
}
