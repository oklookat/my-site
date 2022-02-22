package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"servus/core/internal/iHTTP"
	"servus/core/internal/stacktracer"
	"servus/core/internal/zipify"
)

// get request params.
type httpParamsGetter func(r *http.Request) map[string]string

type httpHelper struct {
	logger  Logger
	control Controller
	cookie  *iHTTP.ConfigCookie
}

func (h *httpHelper) new(l Logger, c Controller, cookie *iHTTP.ConfigCookie, vars httpParamsGetter) {
	iHTTP.RouteArgsGetter = (func(r *http.Request) map[string]string)(vars)
	h.logger = l
	h.control = c
	h.cookie = cookie
}

func (h *httpHelper) getInstance(req *http.Request, res http.ResponseWriter) *iHTTP.Instance {
	var _http *iHTTP.Instance
	var notifyAboutError = func(code int, err error) {
		if err == nil {
			err = errors.New("*empty error*")
		}
		// log.
		h.logger.Error("[core/http] error: " + err.Error())
		// set trace.
		var trace = stacktracer.Instance{}
		trace.Set(err)

		// create zip.
		var zip = zipify.New()
		err1 := zip.AddFile("stacktrace.txt", trace.GetReader())
		err2 := zip.AddFile("request.txt", _http.GetDump())
		// log zip errors.
		if err1 != nil || err2 != nil {
			var message = fmt.Sprintf("[#dump #error%v] http error. Make dump also failed.", code)
			if err1 != nil {
				h.logger.Error("[core/http]: " + err1.Error())
			}
			if err2 != nil {
				h.logger.Error("[core/http]: " + err2.Error())
			}
			h.control.SendMessage(message)
			return
		}
		// send dump.
		var message = fmt.Sprintf("[#dump #error%v] http error.", code)
		h.control.SendFile(&message, trace.GetTimestamp()+".zip", zip.GetRAW())
	}
	// when http error.
	var onHTTPError = func(code int, err error) {
		if code == 500 {
			notifyAboutError(code, err)
		}
	}
	// when response sending error.
	var onSendError = func(code int, err error) {
		h.logger.Error("[core/http]: failed to send response. Error:" + err.Error())
		notifyAboutError(code, err)
	}
	// create iHTTP.
	_http = iHTTP.New(req, res, h.cookie)
	_http.OnHTTPError(onHTTPError)
	_http.OnSendError(onSendError)
	return _http
}

func (h *httpHelper) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var _http = h.getInstance(request, response)
		var ctx = context.WithValue(request.Context(), ctxHTTP, _http)
		*request = *request.WithContext(ctx)
		next.ServeHTTP(response, request)
	})
}
