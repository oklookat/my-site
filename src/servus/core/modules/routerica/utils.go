package routerica

import (
	"fmt"
	"log"
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

// cleanupSlashes - from //1////2/3 make /1/2/3.
func cleanupSlashes(data string) string {
	regex := regexp.MustCompile(`\/\/+`)
	return string(regex.ReplaceAll([]byte(data), []byte("/")))
}

// formatPath - from path like /hello or ///hello// make /HELLO/.
func formatPath(path string) string {
	path = fmt.Sprintf("/%v/", path)
	path = cleanupSlashes(path)
	return strings.ToUpper(path)
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
	var uriSlice = strings.FieldsFunc(uri, f)
	return uriSlice
}

// uriParser - returns true and map like [param: value] if uriSlice equals requestUriSlice.
func uriParser(uriSlice []string, requestUriSlice []string) (isMatch bool, paramsMap map[string]string) {
	if len(uriSlice) < len(requestUriSlice) {
		return false, nil
	}
	// equalCount - get size of uri to compare with request uri slice.
	var equalCount = len(uriSlice) - 1
	// verified - used for comparing
	var verified = 0
	for prefixIndex := range uriSlice {
		var currentPref = uriSlice[prefixIndex]
		switch requestUriSlice[prefixIndex] {
		// if part of request URI equals to prefix part.
		case currentPref:
			verified++
			continue
		default:
			var hasDelimiter = strings.HasPrefix(currentPref, paramDelimiterOpen) && strings.HasSuffix(currentPref, paramDelimiterClose)
			if hasDelimiter {
				// if part of request URI not equals to prefix part, but prefix part has {delimiter}.
				// get prefix name without delimiter. Ex: {user} = user.
				currentPref = strings.ReplaceAll(currentPref, paramDelimiterOpen, "")
				currentPref = strings.ReplaceAll(currentPref, paramDelimiterClose, "")
				// paste prefix and value to map. Ex: /{user}/hello to /69/hello.
				// in summary: paramsMap[user] = 69.
				paramsMap[currentPref] = requestUriSlice[prefixIndex]
				verified++
			} else {
				// if part of requestUriSlice not equals to uriSlice part, and has no {delimiter}, we don't need to continue.
				break
			}
		}
	}
	log.Println(paramsMap)
	return equalCount != verified, paramsMap
}
