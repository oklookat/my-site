package elven

import (
	"net/http"
	"servus/core"
	"servus/core/modules/routerica"
)

func bootRoutes(){
	router := routerica.New()
	router.Use(core.Middleware.MiddlewareCORS, core.Middleware.MiddlewareSecurity, core.Middleware.MiddlewareAsJSON)
	//
	var elvenAuth = router.Group("/elven/auth")
	elvenAuth.POST("/login", controllerAuthLogin)
	elvenAuth.POST("/logout", controllerAuthLogout)
	//
	var elvenArticles = router.Group("/elven/articles")
	elvenArticles.Use(middlewareReadOnly)
	elvenArticles.GET("", controllerArticlesGet)
	elvenArticles.POST("", controllerArticlesPost)
	elvenArticles.PUT("{id}", controllerArticlesPut)
	elvenArticles.DELETE("{id}", controllerArticlesDelete)
	//
	var elvenFiles = router.Group("/elven/files")
	elvenFiles.Use(middlewareAdminOnly)
	elvenFiles.GET("", controllerFilesGet)
	elvenFiles.POST("", controllerFilesPost)
	elvenFiles.PUT("{id}", controllerFilesPut)
	elvenFiles.DELETE("{id}", controllerFilesDelete)
	http.Handle("/", router)
}
