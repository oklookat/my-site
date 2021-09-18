package core

import "net/http"

// HttpResponse wrap response for cool features
type HttpResponse struct {
	http.ResponseWriter
}

func (r HttpResponse) Send(data string, statusCode int) (n int, err error) {
	r.ResponseWriter.WriteHeader(statusCode)
	n, err = r.ResponseWriter.Write([]byte(data))
	if err != nil {
		Logger.Error(err.Error())
	}
	return
}
