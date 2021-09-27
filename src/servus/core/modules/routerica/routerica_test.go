package routerica

import (
	"fmt"
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
	println("/user/id handler")
	response.WriteHeader(200)
	var params = request.Context().Value(ctxPathParams).(map[string]string)
	var formatted = fmt.Sprintf("/user/id handler | id of user: %v", params["id"])
	response.Write([]byte(formatted))
	return
}


func testGroupHandler(response http.ResponseWriter, request *http.Request){
	println("/api/hello handler")
	response.WriteHeader(200)
	response.Write([]byte("/api/hello"))
	return
}