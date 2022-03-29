package way

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	routeHello1 = "/hello"
	routeHello2 = "/hello/world/1/2/3/4/5/"

	routeGroup1     = "/api"
	routeGroup1Full = routeGroup1 + "/" + routeHello1
)

// TODO: add more test & test refactoring

func TestBasicRouting(t *testing.T) {

	var expectedResponse = map[int]string{
		1: "HANDLER 1",
		2: "HANDLER 2",
	}

	var singleMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			t.Log("-- MIDDLEWARE HERE --")
			next.ServeHTTP(response, request)
			return
		})
	}

	var root = New()

	// single routes.
	root.Route("///"+routeHello1+"////", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse[1])
	})
	root.Route(routeHello2+"//", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse[2])
	}).Use(singleMiddleware)

	// init server.
	var requestor = Requestor{}
	requestor.New(root)
	defer requestor.Server.Close()

	// GET 1
	t.Log("// TEST GET 1 //")
	var body, err = requestor.PrettySender(http.MethodGet, routeHello1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse[1] {
		t.Fatalf("expected body: %v, got: %v", expectedResponse[1], body)
	}

	// GET 2
	t.Log("// TEST GET 2 //")
	body, err = requestor.PrettySender(http.MethodGet, routeHello2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse[2] {
		t.Fatalf("expected body: %v, got: %v", expectedResponse[2], body)
	}
}

func TestTypicalRouting(t *testing.T) {
	var rootMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			t.Log("-- ROOT MIDDLEWARE --")
			next.ServeHTTP(response, request)
			return
		})
	}
	var rootMiddleware2 = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			t.Log("-- ROOT MIDDLEWARE 2 --")
			next.ServeHTTP(response, request)
			return
		})
	}
	var logoutMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			t.Log("-- LOGOUT MIDDLEWARE --")
			next.ServeHTTP(response, request)
			return
		})
	}

	var router = New()
	var root = router.Group("/api")
	root.Use(rootMiddleware, rootMiddleware2)

	// auth.
	var authGroup = root.Group("/auth")

	// login
	authGroup.Route("/login", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("/login route")
	}).Methods(http.MethodPost)

	// logout
	var logout = authGroup.Route("/logout", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("/logout route")
	}).Methods(http.MethodPost)
	logout.Use(logoutMiddleware)

	// init server.
	var requestor = Requestor{}
	requestor.New(router)
	defer requestor.Server.Close()

	// REQ 1
	t.Log("// AUTH / LOGIN //")
	var body, err = requestor.PrettySender(http.MethodPost, "/api/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(body)
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
	return r.buildNotSafeRequest(http.MethodPut, path, body)
}

func (r *Requestor) DELETE(path string, body io.Reader) (*http.Response, error) {
	return r.buildNotSafeRequest(http.MethodDelete, path, body)
}

func (r *Requestor) PATCH(path string, body io.Reader) (*http.Response, error) {
	return r.buildNotSafeRequest(http.MethodPatch, path, body)
}

func (r *Requestor) buildNotSafeRequest(method string, path string, body io.Reader) (*http.Response, error) {
	var req, err = http.NewRequest(method, r.URL+path, body)
	if err != nil {
		return nil, err
	}
	res, err := r.Client.Do(req)
	return res, err
}
