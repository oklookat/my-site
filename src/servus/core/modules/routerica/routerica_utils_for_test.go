package routerica

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func readBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("response read failed: %v", err.Error())
	}
	return string(body), nil
}

func Expected(expected string, got string) error {
	if !strings.EqualFold(expected, got){
		return fmt.Errorf("expected: %v, got: %v", expected, got)
	}
	return nil
}
