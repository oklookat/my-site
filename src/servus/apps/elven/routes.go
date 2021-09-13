package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/apps/elven/elUser"
)

func bootRoutes(){
	// TODO: fix MultipartForm and add many variants to get username and password
	// TODO: fix errorCollector casing to normalCasing (get reference from js error collector backend)
	// TODO: test CORS and test cookie max age
	router := mux.NewRouter()
	router.Use(servus.Middleware.MiddlewareAsJSON)
	routerSub := router.PathPrefix("/elven/").Subrouter()
	routerSub.HandleFunc("/auth/login", elUser.ControllerAuthLogin).Methods("POST")
	routerSub.HandleFunc("/auth/logout", elUser.ControllerAuthLogout).Methods("POST")
	//routerSub.HandleFunc("/articles", elControllers.Articles)
	//routerSub.HandleFunc("/articles", elControllers.Files)
	var useCors = servus.Middleware.MiddlewareCORS(router)
	http.Handle("/", useCors)
	//http.Handle("/", router)
}
