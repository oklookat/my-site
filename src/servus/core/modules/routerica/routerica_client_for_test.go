package routerica

import (
	"bytes"
	"net/http"
)


type testI interface {
	newRequest(method string) *http.Request
	TestGET(path string) (response string, err error)
	TestPUT(path string, body string) (responseBody string, err error)
	TestPOST(path string, body string) (responseBody string, err error)
	TestDELETE(path string, body string) (responseBody string, err error)
}

type TestClient struct {
	testI
	url string
	sender *http.Client
}

func NewTestClient(serverURL string) *TestClient {
	return &TestClient{url: serverURL, sender: &http.Client{}}
}

func (c *TestClient) newRequest(method string, path string, body string) *http.Request {
	var reader = bytes.NewReader([]byte(body))
	var request *http.Request
	if method == http.MethodGet {
		request, _ = http.NewRequest(method, c.url + path, nil)
	} else {
		request, _ = http.NewRequest(method, c.url + path, reader)
	}
	return request
}


func (c *TestClient) TestGET(path string) (responseBody string, err error) {
	request := c.newRequest(http.MethodGet, path, "")
	response, err := c.sender.Do(request)
	if err != nil {
		return "", err
	}
	responseBody, err = readBody(response)
	return responseBody, err
}

func (c *TestClient) TestPOST(path string, body string) (responseBody string, err error) {
	request := c.newRequest(http.MethodPost, path, body)
	response, err := c.sender.Do(request)
	if err != nil {
		return "", err
	}
	responseBody, err = readBody(response)
	return responseBody, err
}

func (c *TestClient) TestPUT(path string, body string) (responseBody string, err error) {
	request := c.newRequest(http.MethodPut, path, body)
	response, err := c.sender.Do(request)
	if err != nil {
		return "", err
	}
	responseBody, err = readBody(response)
	return responseBody, err
}

func (c *TestClient) TestDELETE(path string, body string) (responseBody string, err error) {
	request := c.newRequest(http.MethodDelete, path, body)
	response, err := c.sender.Do(request)
	if err != nil {
		return "", err
	}
	responseBody, err = readBody(response)
	return responseBody, err
}
