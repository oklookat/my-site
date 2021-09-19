package core

import (
	"net/http"
)

// HttpResponse wrap response for cool features
type HttpResponse struct {
	http.ResponseWriter
}

func (r HttpResponse) Send(data string, statusCode int){
	r.ResponseWriter.WriteHeader(statusCode)
	_, err := r.ResponseWriter.Write([]byte(data))
	if err != nil {
		Logger.Error(err.Error())
	}
	return
}