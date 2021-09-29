package routerica

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func testRequestBasic(t *testing.T, method string) {
	routerica := New()
	switch method {
	case http.MethodGet:
		routerica.GET(Path1, TestingEndpoint)
		routerica.GET(Path2, TestingEndpoint)
		routerica.GET(Path3, TestingEndpoint)
		break
	case http.MethodPost:
		routerica.POST(Path1, TestingEndpoint)
		routerica.POST(Path2, TestingEndpoint)
		routerica.POST(Path3, TestingEndpoint)
		break
	case http.MethodPut:
		routerica.PUT(Path1, TestingEndpoint)
		routerica.PUT(Path2, TestingEndpoint)
		routerica.PUT(Path3, TestingEndpoint)
		break
	case http.MethodDelete:
		routerica.DELETE(Path1, TestingEndpoint)
		routerica.DELETE(Path2, TestingEndpoint)
		routerica.DELETE(Path3, TestingEndpoint)
		break
	default:
		t.Fatalf("testRequestBasic wrong http method: %v", method)
	}
	server := httptest.NewServer(routerica)
	defer server.Close()
	client := NewTestClient(server.URL)
	var pathFor = func(method string, pathsAndExpected map[string]string) {
		var responseBody string
		var err error
		for path := range pathsAndExpected {
			t.Logf("%v: %v", method, path)
			switch method {
			case http.MethodGet:
				responseBody, err = client.TestGET(path)
				break
			case http.MethodPost:
				responseBody, err = client.TestPOST(path, "")
				break
			case http.MethodPut:
				responseBody, err = client.TestPUT(path, "")
				break
			case http.MethodDelete:
				responseBody, err = client.TestDELETE(path, "")
				break
			}
			var expected = pathsAndExpected[path]
			if err != nil {
				t.Fatalf("routerica %v #1 failed: %v", method, err.Error())
			}
			err = Expected(expected, responseBody)
			if err != nil {
				t.Fatalf(err.Error())
			}
		}
	}
	switch method {
	case http.MethodGet:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response1,
			Path2: Path2Response1,
			Path3: Path3Response1,
		}
		pathFor(http.MethodGet, pathsAndExpected)
		break
	case http.MethodPost:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response2,
			Path2: Path2Response2,
			Path3: Path3Response2,
		}
		pathFor(http.MethodPost, pathsAndExpected)
		break
	case http.MethodPut:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response3,
			Path2: Path2Response3,
			Path3: Path3Response3,
		}
		pathFor(http.MethodPut, pathsAndExpected)
		break
	case http.MethodDelete:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response4,
			Path2: Path2Response4,
			Path3: Path3Response4,
		}
		pathFor(http.MethodDelete, pathsAndExpected)
		break
	}
}

func TestNew(t *testing.T) {
	routerica := New()
	if routerica == nil {
		t.Fatal("Failed to create new instance.")
	}
}

func TestRouterica_GET(t *testing.T) {
	testRequestBasic(t, http.MethodGet)
}

func TestRouterica_POST(t *testing.T) {
	testRequestBasic(t, http.MethodPost)
}

func TestRouterica_PUT(t *testing.T) {
	testRequestBasic(t, http.MethodPut)
}

func TestRouterica_DELETE(t *testing.T) {
	testRequestBasic(t, http.MethodDelete)
}

func TestRouteGroup_GET(t *testing.T) {
	routerica := New()
	server := httptest.NewServer(routerica)
	defer server.Close()
}

func TestRouteGroup_POST(t *testing.T) {
	routerica := New()
	server := httptest.NewServer(routerica)
	defer server.Close()
}

func TestRouteGroup_PUT(t *testing.T) {
	routerica := New()
	server := httptest.NewServer(routerica)
	defer server.Close()
}

func TestRouteGroup_DELETE(t *testing.T) {
	routerica := New()
	server := httptest.NewServer(routerica)
	defer server.Close()
}
