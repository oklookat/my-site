package _app

import (
	"net/http"
)

func bootRoutes(){
	http.HandleFunc("/api/elven/auth/login", auth)
	http.HandleFunc("/api/elven/auth/logout", auth)
	http.HandleFunc("/api/elven/articles", articles)
	http.HandleFunc("/api/elven/files", files)
}
