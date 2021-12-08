package way

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

// defaultNotFound - default endpoint for 404 page.
func defaultNotFound(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
	_, err := response.Write([]byte("not found"))
	if err != nil {
		log.Printf("way: default 404 endpoint, response send failed. Error: %v", err)
		return
	}
	return
}

// defaultNotAllowed - default endpoint for 405 page.
func defaultNotAllowed(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(405)
	_, err := response.Write([]byte("method not allowed"))
	if err != nil {
		log.Printf("way: default 405 endpoint, response send failed. Error: %v", err)
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

// normalizePath - from path like /hello or ///hello// make /hello/.
func normalizePath(path string) string {
	regex := regexp.MustCompile(`//+`)
	path = regex.ReplaceAllString(path, "/")
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
	for index := range maps {
		var _map = maps[index]
		if _map == nil || len(_map) < 1 {
			continue
		}
		for key := range _map {
			concat[key] = _map[key]
		}
	}
	return concat
}

// isEmpty - v len < 1.
func isEmpty(v []string) bool {
	return len(v) < 1
}

// getParamName - check is str is param, and get param name.
func getParamName(str string) (hasParam bool, name string) {
	hasParam = strings.HasPrefix(str, paramOpen) && strings.HasSuffix(str, paramClose)
	if !hasParam {
		return
	}
	name = str
	// remove first {
	name = trimFirstRune(name)
	// remove last }
	name = name[:len(name)-1]
	return
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
