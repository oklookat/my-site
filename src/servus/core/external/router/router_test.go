package router

import (
	"fmt"
	"net/http"
	"testing"
)

func TestManual(t *testing.T) {
	router := New()
	// basic test
	router.Use(testingGlobalMiddleware, testingGlobalMiddleware2)
	var userGroup = router.Group("/api/user")
	userGroup.Use(testingGroupGlobalMiddleware)
	userGroup.GET("{id}", testHandlerGET)
	userGroup.POST("{id}", testHandlerPOST)
	userGroup.PUT("{id}", testHandlerPUT).Use(testingMiddleware)
	userGroup.DELETE("{id}", testHandlerDELETE)
	http.Handle("/", router)
	// http://localhost:7777
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}

func testingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("local middleware 1")
		next.ServeHTTP(writer, request)
	})
}

func testingGlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("global middleware 1")
		next.ServeHTTP(writer, request)
	})
}

func testingGlobalMiddleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("global middleware 2")
		next.ServeHTTP(writer, request)
	})
}

func testingGroupGlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("group middleware 1")
		next.ServeHTTP(writer, request)
	})
}

func testHandlerGET(response http.ResponseWriter, request *http.Request) {
	println("/user/id handler")
	response.WriteHeader(200)
	var params = GetParams(request)
	var formatted = fmt.Sprintf("/user/id handler GET | id of user: %v", params["id"])
	response.Write([]byte(formatted))
	return
}

func testHandlerPOST(response http.ResponseWriter, request *http.Request) {
	println("/user/id handler")
	response.WriteHeader(200)
	var params = GetParams(request)
	var formatted = fmt.Sprintf("/user/id handler POST | id of user: %v", params["id"])
	response.Write([]byte(formatted))
	return
}

func testHandlerPUT(response http.ResponseWriter, request *http.Request) {
	println("/user/id handler")
	response.WriteHeader(200)
	var params = GetParams(request)
	var formatted = fmt.Sprintf("/user/id handler PUT | id of user: %v", params["id"])
	response.Write([]byte(formatted))
	return
}

func testHandlerDELETE(response http.ResponseWriter, request *http.Request) {
	println("/user/id handler")
	response.WriteHeader(200)
	var params = GetParams(request)
	var formatted = fmt.Sprintf("/user/id handler DELETE | id of user: %v", params["id"])
	response.Write([]byte(formatted))
	return
}
