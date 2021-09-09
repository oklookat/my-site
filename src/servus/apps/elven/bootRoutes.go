package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/apps/elven/elControllers"
	"servus/apps/elven/elMiddleware"
)

func bootRoutes(){
	elControllers.BootControllers(servus)
	router := mux.NewRouter()
	router.Use(elMiddleware.AsJSON)
	routerSub := router.PathPrefix("/elven/").Subrouter()
	routerSub.HandleFunc("/auth/login", elControllers.Login).Methods("POST")
	routerSub.HandleFunc("/auth/logout", elControllers.Logout).Methods("POST")
	routerSub.HandleFunc("/articles", elControllers.Articles)
	routerSub.HandleFunc("/articles", elControllers.Files)
	http.Handle("/", router)
}
