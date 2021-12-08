package way

import "net/http"

type route struct {
	path   string
	points []*Point
}

// 	response.Header().Add("Allow", strings.Join(p.methods, ", "))
//	handler405.ServeHTTP(response, request)

func (r *route) new(path string, handler http.HandlerFunc, routes []*route) (notExists bool, p *Point) {
	path = normalizePath(path)
	var point = &Point{}
	point.new(handler)
	// find route.
	var routeExists = false
	for _, rRoute := range routes {
		if rRoute.path != path {
			continue
		}
		routeExists = true
		rRoute.points = append(rRoute.points, point)
	}
	if !routeExists {
		r.path = path
		r.points = make([]*Point, 0)
		r.points = append(r.points, point)
		return true, point
	}
	return false, point
}

// match - this route? If it is, get params.
func (r *route) match(requestPath string) (matched bool, params map[string]string) {
	matched = false
	invalid, differ, params := verifyPaths(r.path, requestPath, true)
	if invalid || differ {
		return
	}
	matched = true
	return
}

func (r *route) findAndExecute(response http.ResponseWriter, request *http.Request, handler405 http.HandlerFunc) {
 	var rMethod = request.Method
	for index := range r.points {
		for _, method := range r.points[index].methods {
			if method != rMethod {
				continue
			}
			r.points[index].runMiddlewareAndEndpoint(response, request)
			return
		}
	}
	handler405.ServeHTTP(response, request)
}
