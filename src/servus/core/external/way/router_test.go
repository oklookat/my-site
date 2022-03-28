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
		3: "GROUP HANDLER 1",
	}

	var singleMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			t.Log("-- MIDDLEWARE HERE --")
			next.ServeHTTP(response, request)
			return
		})
	}

	var groupMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			t.Log("-- GROUP MIDDLEWARE HERE --")
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

	// groups.
	var group = root.Group(routeGroup1)
	group.Use(groupMiddleware)
	group.Route(routeHello1, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResponse[3])
	})

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

	// GET GROUP 1
	t.Log("// TEST GET GROUP 1 //")
	body, err = requestor.PrettySender(http.MethodGet, routeGroup1Full, nil)
	if err != nil {
		t.Fatal(err)
	}
	if body != expectedResponse[3] {
		t.Fatalf("expected body: %v, got: %v", expectedResponse[3], body)
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
