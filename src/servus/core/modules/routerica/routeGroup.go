package routerica

import (
	"fmt"
	"net/http"
	"strings"
)


type routeGroupI interface {
	// Use - add global middleware to route group
	Use(middlewares ...Middleware)
}

type RouteGroup struct {
	RouteBase
	requestsI
	routeGroupI
	prefix      []string
	localRoutes map[string][]RouteLocal
}

// parseRouteGroup - get route group from request URI.
func parseRouteGroup(groups []RouteGroup, requestURI []string) *RouteGroup {
	for index := range groups {
		var currentGroup = &groups[index]
		var isGroupMatched, _ = uriParser(currentGroup.prefix, requestURI)
		if !isGroupMatched {
			continue
		}
		return currentGroup
	}
	return nil
}

// Use - add global middlewares. Any request on this route group will call these middlewares.
func (g *RouteGroup) Use(middlewares ...Middleware) *RouteGroup{
	g.middlewareChain = middlewareChainer(middlewares, &g.RouteBase)
	return g
}

// routeGroupAddRequestHook - executes when adding request to the route group.
func routeGroupAddRequestHook(g *RouteGroup, path *string){
	var prefix = strings.Join(g.prefix, "/")
	*path = fmt.Sprintf("/%v/%v/", prefix, *path)
	*path = formatPath(*path)
}

// GET - add GET request to the route group.
func (g *RouteGroup) GET(path string, handler func(http.ResponseWriter, *http.Request)) *RouteLocal {
	routeGroupAddRequestHook(g, &path)
	var theRoute = newRouteLocal(path, methodGET, handler)
	// add new route to GET group
	g.localRoutes[methodGET] = append(g.localRoutes[methodGET], theRoute)
	var lastRoute = len(g.localRoutes[methodGET]) - 1
	return &g.localRoutes[methodGET][lastRoute]
}
