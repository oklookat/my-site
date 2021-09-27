package routerica

import (
	"log"
	"net/http"
)

// defaultsEndpointNotFound - default endpoint for 404 page.
func defaultsEndpointNotFound(response http.ResponseWriter, request *http.Request){
	response.WriteHeader(404)
	_, err := response.Write([]byte("not found"))
	if err != nil {
		log.Printf("Routerica: default 404 endpoint response send failed. Error: %v", err)
		return
	}
	return
}
