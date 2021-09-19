package routerica

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	routerica := New()
	routerica.Get("/test", testHandler).Use(testingMiddleware, testingMiddleware2)
	http.Handle("/", routerica)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}

func testingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("i am in first middleware /test")
		next.ServeHTTP(writer, request)
	})
}

func testingMiddleware2(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("i am in second middleware /test")
		next.ServeHTTP(writer, request)
	})
}

func testHandler(response http.ResponseWriter, request *http.Request){
	println("i am in handler /test")
	response.WriteHeader(200)
	response.Write([]byte("testing"))
	return
}