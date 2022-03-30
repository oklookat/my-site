package goway

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// TODO: coverage 100

func TestRouting_BasicRootMiddleware(t *testing.T) {
	var executedMiddlewares = 0
	var rootMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			executedMiddlewares++
			next.ServeHTTP(response, request)
			return
		})
	}
	var rootMiddleware2 = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			executedMiddlewares++
			next.ServeHTTP(response, request)
			return
		})
	}
	var root = New()
	root.Use(rootMiddleware, rootMiddleware2)
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	//
	var routePath = "////hello//world//ok"
	var requestPath = "/hello/world/ok"
	//

	var expectedResponse1 = "HANDLER 1"
	root.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse1)
	})

	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse1 {
		t.Fatalf("expected body: %v, got: %v", expectedResponse1, body)
	}
	if executedMiddlewares != 2 {
		t.Fatalf("root middlewares not executed")
	}
}

func TestRouting_RouteMiddleware(t *testing.T) {
	var root = New()
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	var isMiddlewareExecuted = false
	var expectedResponse = "HANDLER 2"
	var routeMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			isMiddlewareExecuted = true
			next.ServeHTTP(response, request)
			return
		})
	}

	//
	var routePath = "i/like/dogs"
	var requestPath = "/i/like/dogs"
	//

	root.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse)
	}).Use(routeMiddleware)

	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse {
		t.Fatalf("expected body: %v, got: %v", expectedResponse, body)
	}
	if !isMiddlewareExecuted {
		t.Fatalf("middleware not executed (GET #2)")
	}
}

func TestRouting_RouteResponseMiddleware(t *testing.T) {
	var root = New()
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	var expectedResponse = "MIDDLEWARE RESPONSE"
	var isHandlerExecuted = false
	var responseSendMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.Write([]byte(expectedResponse))
			return
		})
	}

	//
	var routePath = "/api///users"
	var requestPath = "/api/users"
	//

	root.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		isHandlerExecuted = true
	}).Use(responseSendMiddleware)

	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse {
		t.Fatalf("expected body: %v, got: %v", expectedResponse, body)
	}
	if isHandlerExecuted {
		t.Fatalf("handler executed, but shouldn't")
	}
}

func TestRouting_RootResponseMiddleware(t *testing.T) {
	var expectedResponse = "MIDDLEWARE RESPONSE"
	var responseSendMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			response.Write([]byte(expectedResponse))
			return
		})
	}
	var root = New()
	root.Use(responseSendMiddleware)
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	// test.
	var isHandlerExecuted = false

	//
	var routePath = "another/route/path/ok"
	var requestPath = "another/route/path/ok"
	//

	root.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		isHandlerExecuted = true
	}).Use(responseSendMiddleware)

	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse {
		t.Fatalf("expected body: %v, got: %v", expectedResponse, body)
	}
	if isHandlerExecuted {
		t.Fatalf("handler executed, but shouldn't")
	}
}

func TestRouting_RouteAllowedMethods(t *testing.T) {

	var expectedResponse = map[int]string{
		1: "GET METHOD ROUTE",
		2: "DELETE METHOD ROUTE",
	}

	var root = New()

	// init server.
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	//
	var routePath = "/london//is///the/capital/of//great/britain///"
	var requestPath = "/london/is/the/capital/of/great/britain"
	//

	// GET.
	root.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse[1])
	}).Methods(http.MethodGet)
	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse[1] {
		t.Fatalf("expected body: %v, got: %v", expectedResponse[1], body)
	}

	// DELETE.
	root.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse[2])
	}).Methods(http.MethodDelete)
	body, err = requestor.PrettySender(http.MethodDelete, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse[2] {
		t.Fatalf("expected body: %v, got: %v", expectedResponse[2], body)
	}
}

func TestRouting_GroupAllowedMethods(t *testing.T) {
	var root = New()

	// init server.
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	//
	var groupPath = "/api"
	var routePath = "/users/add"
	var requestPath = "/api/users/add"
	//

	var group = root.Group(groupPath).Methods(http.MethodGet, http.MethodDelete)

	var expectedResponse1 = "HANDLER"
	group.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse1)
	})

	// allowed method (GET).
	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse1 {
		t.Fatalf("expected body: %v, got: %v", expectedResponse1, body)
	}

	// allowed method (DELETE).
	body, err = requestor.PrettySender(http.MethodDelete, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse1 {
		t.Fatalf("expected body: %v, got: %v", expectedResponse1, body)
	}

	// Not allowed method.
	var notAllowedExpected = "method not allowed"
	Handler405 = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(notAllowedExpected))
	}
	body, err = requestor.PrettySender(http.MethodPost, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != notAllowedExpected {
		t.Fatalf("expected body: %v, got: %v", notAllowedExpected, body)
	}

}

func TestRouting_EmptyRoutesPathInGroup(t *testing.T) {
	var root = New()

	// init server.
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	//
	var groupPath = "group/path"
	var groupPath2 = "second/group/"

	var requestPath2 = "/group/path/second/group/"
	//

	var group = root.Group(groupPath)
	var groupUnderGroup = group.Group(groupPath2)

	// PATCH.
	var expectedResponse1 = "PATCH METHOD ROUTE"
	groupUnderGroup.Route("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse1)
	}).Methods(http.MethodPatch)
	var body, err = requestor.PrettySender(http.MethodPatch, requestPath2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse1 {
		t.Fatalf("expected body: %v, got: %v", expectedResponse1, body)
	}

	// GET.
	var expectedResponse2 = "GET METHOD ROUTE"
	groupUnderGroup.Route("", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse2)
	}).Methods(http.MethodGet)
	body, err = requestor.PrettySender(http.MethodGet, requestPath2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse2 {
		t.Fatalf("expected body: %v, got: %v", expectedResponse2, body)
	}

	// Not allowed method.
	var notAllowedExpected = "method not allowed"
	Handler405 = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(notAllowedExpected))
	}
	body, err = requestor.PrettySender(http.MethodPost, requestPath2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != notAllowedExpected {
		t.Fatalf("expected body: %v, got: %v", notAllowedExpected, body)
	}
}

func TestRouting_Router404(t *testing.T) {
	var root = New()

	// init server.
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	//
	var requestPath = "/hello/world/1234"
	var expectedResponse = "404 handler"
	//

	Handler404 = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedResponse))
	}

	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse {
		t.Fatalf("expected body: %v, got: %v", Handler404, body)
	}
}

// group vars + route vars.
func TestRouting_Vars(t *testing.T) {
	var root = New()

	// init server.
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	//
	var groupPath = "{username}/actions/{id}/"
	var routePath = "{word}/hello/{who}"
	var requestPath = "/oklookat/actions/69/coolword/hello/world"
	//

	var expectedResponse1 = "OK DONE"
	var expectedVars = map[string]string{
		"username": "oklookat",
		"id":       "69",
		"word":     "coolword",
		"who":      "world",
	}
	var group = root.Group(groupPath)

	group.Route(routePath, func(w http.ResponseWriter, r *http.Request) {
		var rqVars = Vars(r)
		var isSame = reflect.DeepEqual(rqVars, expectedVars)
		if !isSame {
			t.Fatalf("request vars not same")
		}
		fmt.Fprint(w, expectedResponse1)
	})

	var body, err = requestor.PrettySender(http.MethodGet, requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse1 {
		t.Fatalf("expected body: %v, got: %v", expectedResponse1, body)
	}
}

////////////////////////////
type Requestor struct {
	Server *httptest.Server
	Client *http.Client
	URL    string
}

func (r *Requestor) New(handler http.Handler) {
	r.Server = httptest.NewServer(handler)
	r.URL = r.Server.URL
	r.Client = r.Server.Client()
}

func (r *Requestor) PrettySender(method string, path string, reqBody io.Reader) (body string, err error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	var res *http.Response
	switch method {
	default:
		err = errors.New("wrong request method")
		return
	case http.MethodGet:
		res, err = r.GET(path)
		break
	case http.MethodHead:
		res, err = r.HEAD(path)
		break
	case http.MethodPost:
		res, err = r.POST(path, reqBody)
		break
	case http.MethodPut:
		res, err = r.PUT(path, reqBody)
		break
	case http.MethodDelete:
		res, err = r.DELETE(path, reqBody)
		break
	case http.MethodPatch:
		res, err = r.PATCH(path, reqBody)
		break
	}
	if err != nil || res == nil {
		return "", err
	}
	var byteBody []byte
	byteBody, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()
	body = string(byteBody)
	return
}

func (r *Requestor) GET(path string) (*http.Response, error) {
	var res, err = r.Client.Get(r.URL + path)
	return res, err
}

func (r *Requestor) HEAD(path string) (*http.Response, error) {
	var res, err = r.Client.Head(r.URL + path)
	return res, err
}

func (r *Requestor) POST(path string, body io.Reader) (*http.Response, error) {
	var res, err = r.Client.Post(r.URL+path, "application/json", body)
	return res, err
}

func (r *Requestor) PUT(path string, body io.Reader) (*http.Response, error) {
	return r.buildRequest(http.MethodPut, path, body)
}

func (r *Requestor) DELETE(path string, body io.Reader) (*http.Response, error) {
	return r.buildRequest(http.MethodDelete, path, body)
}

func (r *Requestor) PATCH(path string, body io.Reader) (*http.Response, error) {
	return r.buildRequest(http.MethodPatch, path, body)
}

func (r *Requestor) buildRequest(method string, path string, body io.Reader) (*http.Response, error) {
	var req, err = http.NewRequest(method, r.URL+path, body)
	if err != nil {
		return nil, err
	}
	res, err := r.Client.Do(req)
	return res, err
}
