package core

import (
	"net/http"
)

// HttpResponse wrap response for cool features
type HttpResponse struct {
	http.ResponseWriter
}

func (r *HttpResponse) Send(body string, statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	_, err := r.ResponseWriter.Write([]byte(body))
	if err != nil {
		Logger.Error(err.Error())
	}
	return
}
