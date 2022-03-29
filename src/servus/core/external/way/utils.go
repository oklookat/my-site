package way

import (
	"context"
	"net/http"
	"path"
	"strings"
)

// when ServeHTTP called, sets Executed to true.
type dummyEndpoint struct {
	Executed bool
}

func (d *dummyEndpoint) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	d.Executed = true
}

// run middlewares.
func executeMiddleware(response http.ResponseWriter, request *http.Request, middleware MiddlewareFunc) (isResponseSended bool) {
	if middleware == nil {
		return
	}

	// set dummy as middleware endpoint.
	var dummy = &dummyEndpoint{}
	var handler = middleware(dummy)
	handler.ServeHTTP(response, request)

	// if dummy not executed
	// it means middleware sended response / not called next.ServeHTTP()
	isResponseSended = !dummy.Executed

	return
}

// https://gist.github.com/husobee/fd23681261a39699ee37?permalink_comment_id=3111569#gistcomment-3111569
//
// make one big middleware from middlewares.
func middlewareChain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	if len(middlewares) < 1 {
		return nil
	}
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] == nil {
				continue
			}
			next = middlewares[i](next)
		}
		return next
	}
}

// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
//
// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// get path variables.
func Vars(request *http.Request) map[string]string {
	var ctx = request.Context()
	var varsMap, ok = ctx.Value(CTX_VARS_NAME).(map[string]string)
	if !ok {
		return nil
	}
	return varsMap
}

// is route variable?
func isRouteVar(path string) (isVar bool, varName string) {
	isVar = strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}")
	if isVar {
		var withoutBrackets = strings.ReplaceAll(path, "{", "")
		withoutBrackets = strings.ReplaceAll(withoutBrackets, "}", "")
		varName = withoutBrackets
	}
	return
}

// add route variable with value to request context.
func addVarToContext(request *http.Request, name string, value string) {
	var oldCtx = request.Context()
	var varsMap, ok = oldCtx.Value(CTX_VARS_NAME).(map[string]string)
	if !ok || varsMap == nil {
		varsMap = make(map[string]string, 0)
	}
	varsMap[name] = value
	var ctxWithVars = context.WithValue(oldCtx, CTX_VARS_NAME, varsMap)
	*request = *request.WithContext(ctxWithVars)
}

// check is method allowed.
func isMethodAllowed(methods []string, requestMethod string) bool {
	if methods == nil {
		return true
	}
	for _, method := range methods {
		if requestMethod == method {
			return true
		}
	}
	return false
}

// remove slash at start and end of str.
func removeSlashStartEnd(str string) string {
	str = removeSlashStart(str)
	str = removeSlashEnd(str)
	return str
}

// remove slash at start of str.
func removeSlashStart(str string) string {
	str = strings.TrimPrefix(str, "/")
	return str
}

// remove slash at end of str.
func removeSlashEnd(str string) string {
	str = strings.TrimSuffix(str, "/")
	return str
}

// default 405 sender.
func send405(r http.ResponseWriter) {
	r.WriteHeader(405)
	r.Write([]byte("method not allowed"))
}

// default 404 sender.
func send404(r http.ResponseWriter) {
	r.WriteHeader(404)
	r.Write([]byte("not found"))
}

// make path like: /hello/world
func pathToStandart(to string) string {
	if len(to) < 1 {
		return to
	}
	to = "/" + removeSlashEnd(to)
	return path.Clean(to)
}

// https://www.geeksforgeeks.org/how-to-remove-duplicate-values-from-slice-in-golang/
//
// remove duplicates from slice.
func removeDuplicateValues[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// format route allowed methods.
func processAllowedMethods(slice []string, methods ...string) []string {
	if slice == nil {
		slice = make([]string, 0)
	}
	for _, method := range methods {
		method = strings.ToUpper(method)
		slice = append(slice, method)
	}
	slice = removeDuplicateValues(slice)
	return slice
}
