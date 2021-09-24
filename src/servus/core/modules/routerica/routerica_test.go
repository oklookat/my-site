package routerica

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	routerica := New()
	// basic test
	routerica.Use(testingGlobalMiddleware, testingGlobalMiddleware2)
	routerica.GET("/test", testHandler).Use(testingMiddleware, testingMiddleware2)
	// group test
	var group = routerica.Group("/api").Use(testGroupGlobalMiddleware)
	group.GET("/hello", testGroupHandler).Use(testGroupMiddleware)
	//
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

func testingGlobalMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("i am in global middleware")
		next.ServeHTTP(writer, request)
	})
}

func testingGlobalMiddleware2(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("i am in global middleware 2")
		next.ServeHTTP(writer, request)
	})
}

func testHandler(response http.ResponseWriter, request *http.Request){
	println("i am in handler /test")
	response.WriteHeader(200)
	response.Write([]byte("testing"))
	return
}

func testGroupGlobalMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("i am in global group middleware /api/")
		next.ServeHTTP(writer, request)
	})
}

func testGroupMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("i am in group middleware /api/hello")
		next.ServeHTTP(writer, request)
	})
}

func testGroupHandler(response http.ResponseWriter, request *http.Request){
	println("i am in group handler /api/hello")
	response.WriteHeader(200)
	response.Write([]byte("testing"))
	return
}