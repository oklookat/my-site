package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
)

func bootRoutes() {
	router := mux.NewRouter()
	router.Use(core.Middleware.MiddlewareAsJSON)
	var routerElven = router.PathPrefix("/elven").Subrouter()
	//
	var routerAuth = routerElven.PathPrefix("/auth").Subrouter()
	routerAuth.HandleFunc("/login", eAuth.controllerLogin).Methods(http.MethodPost)
	routerAuth.HandleFunc("/logout", eAuth.controllerLogout).Methods(http.MethodPost)
	//
	var routerArticles = routerElven.PathPrefix("/articles").Subrouter()
	routerArticles.Use(eBase.middlewareReadOnly)
	routerArticles.HandleFunc("", eArticle.controllerGetAll).Methods(http.MethodGet)
	routerArticles.HandleFunc("/{id}", eArticle.controllerGetOne).Methods(http.MethodGet)
	routerArticles.HandleFunc("", eArticle.controllerCreateOne).Methods(http.MethodPost)
	routerArticles.HandleFunc("/{id}", eArticle.controllerUpdateOne).Methods(http.MethodPut)
	routerArticles.HandleFunc("/{id}", eArticle.controllerDeleteOne).Methods(http.MethodDelete)
	//
	var routerFiles = routerElven.PathPrefix("/files").Subrouter()
	routerFiles.Use(eBase.middlewareAdminOnly)
	routerFiles.HandleFunc("", eFile.controllerGetAll).Methods(http.MethodGet)
	routerFiles.HandleFunc("/{id}", eFile.controllerGetOne).Methods(http.MethodGet)
	routerFiles.HandleFunc("", eFile.controllerCreateOne).Methods(http.MethodPost)
	routerFiles.HandleFunc("/{id}", eFile.controllerUpdateOne).Methods(http.MethodPut)
	routerFiles.HandleFunc("/{id}", eFile.controllerDeleteOne).Methods(http.MethodDelete)
	//
	var useBeforeRouter = core.Middleware.MiddlewareCORS(core.Middleware.MiddlewareSecurity(router))
	http.Handle("/", useBeforeRouter)
}
