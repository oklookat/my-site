package _app

import (
	"net/http"
	"servus/_app/controllers"
)

func bootRoutes(){
	http.HandleFunc("/api/elven/auth/login", controllers.Auth)
	http.HandleFunc("/api/elven/auth/logout", controllers.Auth)
	http.HandleFunc("/api/elven/articles", controllers.Articles)
	http.HandleFunc("/api/elven/files", controllers.Files)
}
