package routerica

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	routerica := New()
	// basic test
	routerica.Use(testingGlobalMiddleware)
	routerica.GET("/user/{id}", testHandler)
	// group test
	var group = routerica.Group("/api")
	group.GET("/hello", testGroupHandler)
	//
	http.Handle("/", routerica)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}

func testingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("/test first middleware")
		next.ServeHTTP(writer, request)
	})
}


func testingGlobalMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("global middleware 1")
		next.ServeHTTP(writer, request)
	})
}


func testHandler(response http.ResponseWriter, request *http.Request){
	println("/test/get handler")
	response.WriteHeader(200)
	response.Write([]byte("/test"))
	return
}


func testGroupHandler(response http.ResponseWriter, request *http.Request){
	println("/api/hello handler")
	response.WriteHeader(200)
	response.Write([]byte("/api/hello"))
	return
}