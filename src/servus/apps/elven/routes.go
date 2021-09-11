package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

func bootRoutes(){
	router := mux.NewRouter()
	router.Use(middlewareAsJSON)
	routerSub := router.PathPrefix("/elven/").Subrouter()
	routerSub.HandleFunc("/auth/login", controllerAuthLogin).Methods("POST")
	routerSub.HandleFunc("/auth/logout", controllerAuthLogout).Methods("POST")
	//routerSub.HandleFunc("/articles", elControllers.Articles)
	//routerSub.HandleFunc("/articles", elControllers.Files)
	http.Handle("/", router)
}
