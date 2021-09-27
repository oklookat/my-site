package routerica

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// create one big middleware from middlewares (chain).
func middlewareChainer(middlewares []Middleware, next http.Handler) http.Handler {
	// https://gist.github.com/husobee/fd23681261a39699ee37#gistcomment-3111569
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}


// formatPath - from path like /hello or ///hello// make /HELLO/.
func formatPath(path string) string {
	path = fmt.Sprintf("/%v/", path)
	regex := regexp.MustCompile(`//+`)
	path = string(regex.ReplaceAll([]byte(path), []byte("/")))
	return path
}

//func paramMatching(){
//	var regex, _ = regexp.Compile(`{([a-zA-Z0-9]*)}`)
//}

// uriSplitter - split url like /hello/world/ to slice [hello, world].
func uriSplitter(uri string) []string {
	uri = formatPath(uri)
	f := func(c rune) bool {
		return c == '/'
	}
	// split uri by slash.
	var uriSlice = strings.FieldsFunc(uri, f)
	return uriSlice
}

// mapConcat - make one map from maps (duplicates will be replaced)
func mapConcat(maps ...map[string]string) map[string]string{
	mapped := make(map[string]string, 0)
	for currentMapIndex := range maps {
		var currentMap = maps[currentMapIndex]
		for currentMapKey := range currentMap {
			mapped[currentMapKey] = currentMap[currentMapKey]
		}
	}
	return mapped
}
