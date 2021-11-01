package elven

import (
	"github.com/gorilla/mux"
	"net/http"
)

func bootRoutes() {
	router := mux.NewRouter()
	router.Use(instance.HTTP.Middleware.AsJSON)
	var routerElven = router.PathPrefix("/elven").Subrouter()
	//
	var routerAuth = routerElven.PathPrefix("/auth").Subrouter()
	routerAuth.HandleFunc("/login", eAuth.controllerLogin).Methods(http.MethodPost)
	var routerAuthLogout = routerAuth.PathPrefix("/logout").Subrouter()
	routerAuthLogout.Use(eBase.middlewareAuthorizedOnly)
	routerAuthLogout.HandleFunc("", eAuth.controllerLogout).Methods(http.MethodPost)
	//
	var routerArticles = routerElven.PathPrefix("/articles").Subrouter()
	routerArticles.Use(eBase.middlewareReadOnly)
	routerArticles.HandleFunc("", eArticle.getAll).Methods(http.MethodGet)
	routerArticles.HandleFunc("/{id}", eArticle.getOne).Methods(http.MethodGet)
	routerArticles.HandleFunc("", eArticle.create).Methods(http.MethodPost)
	routerArticles.HandleFunc("/{id}", eArticle.update).Methods(http.MethodPut, http.MethodPatch)
	routerArticles.HandleFunc("/{id}", eArticle.delete).Methods(http.MethodDelete)
	//
	var routerFiles = routerElven.PathPrefix("/files").Subrouter()
	routerFiles.Use(eBase.middlewareAdminOnly)
	routerFiles.HandleFunc("", eFile.controllerGetAll).Methods(http.MethodGet)
	routerFiles.HandleFunc("", eFile.controllerCreateOne).Methods(http.MethodPost)
	routerFiles.HandleFunc("/{id}", eFile.controllerDeleteOne).Methods(http.MethodDelete)
	//
	var routerUsers = routerElven.PathPrefix("/users").Subrouter()
	routerUsers.Use(eBase.middlewareAuthorizedOnly)
	routerUsers.HandleFunc("/me", eUser.controllerGetMe).Methods(http.MethodGet)
	routerUsers.HandleFunc("/me/change", eUser.controllerChange).Methods(http.MethodPost)
	//
	var useBeforeRouter = instance.HTTP.Middleware.CORS(instance.HTTP.Middleware.Security(router))
	http.Handle("/", useBeforeRouter)
}

