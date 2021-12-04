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
	var auth = authRoutes{}
	auth.boot(routerElven)
	//
	var article = articleRoutes{}
	article.boot(routerElven)
	//
	var file = fileRoutes{}
	file.boot(routerElven)
	//
	var user = userRoutes{}
	user.boot(routerElven)
	//
	var useBeforeRouter = instance.HTTP.Middleware.CORS(instance.HTTP.Middleware.Security(router))
	http.Handle("/", useBeforeRouter)
}

