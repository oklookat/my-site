package elven

import (
	"github.com/gorilla/mux"
	"net/http"
	"servus/core"
)

func bootRoutes(){
	router := mux.NewRouter()
	router.Use(core.Middleware.MiddlewareAsJSON)
	var elven = router.PathPrefix("/elven").Subrouter()
	//
	var elvenAuth = elven.PathPrefix("/auth").Subrouter()
	elvenAuth.HandleFunc("/login", controllerAuthLogin).Methods(http.MethodPost)
	elvenAuth.HandleFunc("/logout", controllerAuthLogin).Methods(http.MethodPost)
	//
	var elvenArticles = elven.PathPrefix("/articles").Subrouter()
	elvenArticles.Use(middlewareReadOnly)
	elvenArticles.HandleFunc("", controllerArticlesGetAll).Methods(http.MethodGet)
	elvenArticles.HandleFunc("", controllerArticlesPost).Methods(http.MethodPost)
	elvenArticles.HandleFunc("/{id}", controllerArticlesPut).Methods(http.MethodPut)
	elvenArticles.HandleFunc("/{id}", controllerArticlesGetOne).Methods(http.MethodGet)
	elvenArticles.HandleFunc("/{id}", controllerArticlesDelete).Methods(http.MethodDelete)
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
