package routerica

import (
	"net/http"
)

const (
	//
	Path1          = "/hello"
	Path1Response1 = "TPR_1_GET"
	Path1Response2 = "TPR_1_POST"
	Path1Response3 = "TPR_1_PUT"
	Path1Response4 = "TPR_1_DELETE"
	//
	Path2          = "/hello/big/world"
	Path2Response1 = "TPR_2_GET"
	Path2Response2 = "TPR_2_POST"
	Path2Response3 = "TPR_2_PUT"
	Path2Response4 = "TPR_2_DELETE"
	//
	Path3          = "/wow/very/big/world/with/long/paths/and/its/good/or/maybe/not"
	Path3Response1 = "TPR_3_GET"
	Path3Response2 = "TPR_3_POST"
	Path3Response3 = "TPR_3_PUT"
	Path3Response4 = "TPR_3_DELETE"
	//
	TestingPath4 = "/route/with/{username}/and/{password}"
	TestingPath5 = "/{username}/and/{password}"
)

// TestingEndpoint - basic endpoint for test requests.
func TestingEndpoint(response http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case Path1:
		pathHelper(response, request, Path1Response1, Path1Response2, Path1Response3, Path1Response4)
		break
	case Path2:
		pathHelper(response, request, Path2Response1, Path2Response2, Path2Response3, Path2Response4)
		break
	case Path3:
		pathHelper(response, request, Path3Response1, Path3Response2, Path3Response3, Path3Response4)
		break
	}
}

// pathHelper - depending on request method send response (watch TestingEndpoint for example).
func pathHelper(response http.ResponseWriter, request *http.Request, TPR1 string, TPR2 string, TPR3 string, TPR4 string){
	switch request.Method {
	case http.MethodGet:
		response.WriteHeader(200)
		response.Write([]byte(TPR1))
		break
	case http.MethodPost:
		response.WriteHeader(200)
		response.Write([]byte(TPR2))
		break
	case http.MethodPut:
		response.WriteHeader(200)
		response.Write([]byte(TPR3))
		break
	case http.MethodDelete:
		response.WriteHeader(200)
		response.Write([]byte(TPR4))
		break
	}
}

func TestingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("local middleware 1")
		next.ServeHTTP(writer, request)
	})
}

func TestingGlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("global middleware 1")
		next.ServeHTTP(writer, request)
	})
}

func TestingGroupGlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		println("group middleware 1")
		next.ServeHTTP(writer, request)
	})
}
