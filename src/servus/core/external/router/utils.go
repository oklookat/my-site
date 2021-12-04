package router

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// defaultNotFound - default endpoint for 404 page.
func defaultNotFound(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
	_, err := response.Write([]byte("not found"))
	if err != nil {
		log.Printf("router: default 404 endpoint, response send failed. Error: %v", err)
		return
	}
	return
}

// defaultNotAllowed - default endpoint for 405 page.
func defaultNotAllowed(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(405)
	_, err := response.Write([]byte("method not allowed"))
	if err != nil {
		log.Printf("router: default 405 endpoint, response send failed. Error: %v", err)
		return
	}
	return
}

// GetParams - get {params} from request.
//
// Example: if route /hello/{id} and request are /hello/12 - returns [id: 12].
func GetParams(request *http.Request) map[string]string {
	params, ok := request.Context().Value(ctxPathParams).(map[string]string)
	if !ok {
		return nil
	}
	return params
}

// middlewareChainer - create one big middleware from middlewares (chain).
func middlewareChainer(middlewares []Middleware, next http.Handler) http.Handler {
	// https://gist.github.com/husobee/fd23681261a39699ee37#gistcomment-3111569
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

// normalizePath - from path like /hello or ///hello// make /HELLO/.
func normalizePath(path string) string {
	path = fmt.Sprintf("/%v/", path)
	regex := regexp.MustCompile(`//+`)
	path = string(regex.ReplaceAll([]byte(path), []byte("/")))
	return path
}

// pathToSlice - split path like /hello/world/ to slice [hello, world].
func pathToSlice(path string) []string {
	path = normalizePath(path)
	f := func(c rune) bool {
		return c == '/'
	}
	// split uri by slash.
	var uriSlice = strings.FieldsFunc(path, f)
	return uriSlice
}

// mapsToMap - make one map from maps (duplicates will be replaced).
func mapsToMap(maps ...map[string]string) map[string]string {
	concat := make(map[string]string, 0)
	for mapIndex := range maps {
		var _map = maps[mapIndex]
		for mapKey := range _map {
			concat[mapKey] = _map[mapKey]
		}
	}
	return concat
}

// moreThan - if len v1 > v2.
func moreThan(v1 []string, v2 []string) bool {
	return len(v1) > len(v2)
}

// isEmpty - v len < 1.
func isEmpty(v []string) bool {
	return len(v) < 1
}

// addRoute - put route in routes.
func addRoute(routes routes, path string, method string, endpoint http.HandlerFunc) *RouteMethod {
	path = normalizePath(path)
	_route, routeExists := routes[path]
	if !routeExists {
		_route = &route{}
		_route.new()
	}
	var routeMethod = &RouteMethod{}
	routeMethod.new(_route, path, endpoint)
	_route.addMethod(method, routeMethod)
	routes[path] = _route
	return routeMethod
}
