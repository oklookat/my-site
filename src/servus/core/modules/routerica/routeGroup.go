package routerica

import (
	"fmt"
	"net/http"
	"strings"
)


// get route group from request URI
func parseRouteGroup(groups []RouteGroup, requestURI string) *RouteGroup {
	requestURI = strings.ToUpper(requestURI)
	for index := range groups {
		var currentGroup = &groups[index]
		if len(currentGroup.prefix) <= 1 {
			continue
		}
		// if request url starts with current group url
		isRouteGroup := strings.HasPrefix(requestURI+"/", currentGroup.prefix)
		if !isRouteGroup {
			continue
		}
		return currentGroup
	}
	return nil
}

type routeGroupI interface {
	// Use - add global middleware to route group
	Use(middlewares ...Middleware)
}

type RouteGroup struct {
	RouteBase
	requestsI
	routeGroupI
	prefix      string
	localRoutes map[string][]RouteLocal
}

// Use middleware (global group). That method executes when you're adding global middleware to the route group
func (g *RouteGroup) Use(middlewares ...Middleware) *RouteGroup{
	g.middlewareChain = middlewareChainer(middlewares, &g.RouteBase)
	return g
}

func (g *RouteGroup) GET(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal {
	path = fmt.Sprintf("/%v/%v/", g.prefix, path)
	cleanupSlashes(&path)
	var theRoute = RouteLocal{path: path, method: methodGET, RouteBase: RouteBase{handler: handler}}
	g.localRoutes[methodGET] = append(g.localRoutes[methodGET], theRoute)
	var lastRoute = len(g.localRoutes[methodGET]) - 1
	return &g.localRoutes[methodGET][lastRoute]
}
