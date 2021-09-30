package routerica

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func testRequestBasic(t *testing.T, method string) {
	routerica := New()
	var group = routerica.Group(GroupPath1)
	var groupParams = routerica.Group(GroupPath2)
	switch method {
	case http.MethodGet:
		routerica.GET(Path1, TestingEndpoint)
		routerica.GET(Path2, TestingEndpoint)
		routerica.GET(Path3, TestingEndpoint)
		routerica.GET(Path4, TestingEndpoint)
		routerica.GET(Path5, TestingEndpoint)
		group.GET(GroupSubPath1, TestingEndpoint)
		groupParams.GET(GroupSubPath2, TestingEndpoint)
		break
	case http.MethodPost:
		routerica.POST(Path1, TestingEndpoint)
		routerica.POST(Path2, TestingEndpoint)
		routerica.POST(Path3, TestingEndpoint)
		routerica.POST(Path4, TestingEndpoint)
		routerica.POST(Path5, TestingEndpoint)
		group.POST(GroupSubPath1, TestingEndpoint)
		groupParams.POST(GroupSubPath2, TestingEndpoint)
		break
	case http.MethodPut:
		routerica.PUT(Path1, TestingEndpoint)
		routerica.PUT(Path2, TestingEndpoint)
		routerica.PUT(Path3, TestingEndpoint)
		routerica.PUT(Path4, TestingEndpoint)
		routerica.PUT(Path5, TestingEndpoint)
		group.PUT(GroupSubPath1, TestingEndpoint)
		groupParams.PUT(GroupSubPath2, TestingEndpoint)
		break
	case http.MethodDelete:
		routerica.DELETE(Path1, TestingEndpoint)
		routerica.DELETE(Path2, TestingEndpoint)
		routerica.DELETE(Path3, TestingEndpoint)
		routerica.DELETE(Path4, TestingEndpoint)
		routerica.DELETE(Path5, TestingEndpoint)
		group.DELETE(GroupSubPath1, TestingEndpoint)
		groupParams.DELETE(GroupSubPath2, TestingEndpoint)
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
			Path4Params: Path4Response1,
			Path5Params: Path5Response1,
			GroupSubPath1Request: GroupSubPath1Response1,
			GroupSubPath2Request: GroupSubPath2Response1,
		}
		pathFor(http.MethodGet, pathsAndExpected)
		break
	case http.MethodPost:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response2,
			Path2: Path2Response2,
			Path3: Path3Response2,
			Path4Params: Path4Response2,
			Path5Params: Path5Response2,
			GroupSubPath1Request: GroupSubPath1Response2,
			GroupSubPath2Request: GroupSubPath2Response2,
		}
		pathFor(http.MethodPost, pathsAndExpected)
		break
	case http.MethodPut:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response3,
			Path2: Path2Response3,
			Path3: Path3Response3,
			Path4Params: Path4Response3,
			Path5Params: Path5Response3,
			GroupSubPath1Request: GroupSubPath1Response3,
			GroupSubPath2Request: GroupSubPath2Response3,
		}
		pathFor(http.MethodPut, pathsAndExpected)
		break
	case http.MethodDelete:
		var pathsAndExpected = map[string]string{
			Path1: Path1Response4,
			Path2: Path2Response4,
			Path3: Path3Response4,
			Path4Params: Path4Response4,
			Path5Params: Path5Response4,
			GroupSubPath1Request: GroupSubPath1Response4,
			GroupSubPath2Request: GroupSubPath2Response4,
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

func TestRouterica_All_GET(t *testing.T) {
	testRequestBasic(t, http.MethodGet)
}

func TestRouterica_All_POST(t *testing.T) {
	testRequestBasic(t, http.MethodPost)
}

func TestRouterica_All_PUT(t *testing.T) {
	testRequestBasic(t, http.MethodPut)
}

func TestRouterica_All_DELETE(t *testing.T) {
	testRequestBasic(t, http.MethodDelete)
}