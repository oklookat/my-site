package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
)

func bootRoutes(){
	var baseController = &baseController{logger: core.Logger}
	router := mux.NewRouter()
	router.Use(core.Middleware.MiddlewareAsJSON)
	var elven = router.PathPrefix("/elven").Subrouter()
	//
	var authController = controllerAuth{baseController}
	var elvenAuth = elven.PathPrefix("/auth").Subrouter()
	elvenAuth.HandleFunc("/login", authController.login).Methods(http.MethodPost)
	elvenAuth.HandleFunc("/logout", authController.logout).Methods(http.MethodPost)
	//
	var articlesController = controllerArticles{baseController}
	var elvenArticles = elven.PathPrefix("/articles").Subrouter()
	elvenArticles.Use(middlewareReadOnly)
	elvenArticles.HandleFunc("", articlesController.GetAll).Methods(http.MethodGet)
	elvenArticles.HandleFunc("", articlesController.Create).Methods(http.MethodPost)
	elvenArticles.HandleFunc("/{id}", articlesController.Update).Methods(http.MethodPut)
	elvenArticles.HandleFunc("/{id}", articlesController.GetOne).Methods(http.MethodGet)
	elvenArticles.HandleFunc("/{id}", articlesController.Delete).Methods(http.MethodDelete)
	//
	var elvenFiles = elven.PathPrefix("/files").Subrouter()
	elvenFiles.Use(middlewareAdminOnly)
	elvenFiles.HandleFunc("", controllerFilesGet).Methods(http.MethodGet)
	elvenFiles.HandleFunc("", controllerFilesPost).Methods(http.MethodPost)
	elvenFiles.HandleFunc("/{id}", controllerFilesPut).Methods(http.MethodPut)
	elvenFiles.HandleFunc("/{id}", controllerFilesDelete).Methods(http.MethodDelete)
	//
	var useBeforeRouter = core.Middleware.MiddlewareCORS(core.Middleware.MiddlewareSecurity(router))
	http.Handle("/", useBeforeRouter)
}
