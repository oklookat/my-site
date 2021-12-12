package core

import (
	"context"
	"fmt"
	"net/http"
	"servus/core/internal/iHTTP"
	"servus/core/internal/zipify"
)

type HTTP struct {
	logger Logger
	control *Control
	cookie *iHTTP.ConfigCookie
}

func (h *HTTP) new(l Logger, c *Control, cookie *iHTTP.ConfigCookie) {
	h.logger = l
	h.control = c
	h.cookie = cookie
}

func (h *HTTP) getInstance(req *http.Request, res http.ResponseWriter)  *iHTTP.Instance {
	var _http *iHTTP.Instance
	var notifyAboutError = func(code int, err error) {
		// log.
		h.logger.Error("[provide/http] error: " + err.Error())
		// set trace.
		var trace = StackTrace{}
		trace.Set(err)
		// create zip.
		var zip = zipify.New()
		err1 := zip.AddFile("stacktrace.txt", trace.GetReader())
		err2 := zip.AddFile("request.txt", _http.GetDump())
		// log zip errors.
		if err1 != nil || err2 != nil {
			var message = fmt.Sprintf("[#dump #code%v] http error. Make dump also failed.", code)
			if err1 != nil {
				h.logger.Error("[provide/http]: " + err1.Error())
			}
			if err2 != nil {
				h.logger.Error("[provide/http]: " + err2.Error())
			}
			h.control.Telegram.Bot.SendMessage(message)
			return
		}
		// send dump.
		var message = fmt.Sprintf("[#dump #code%v] http error.", code)
		h.control.Telegram.Bot.SendFile(&message, trace.GetTimestamp()+".zip", zip.GetRAW())
	}
	// when http error.
	var onHTTPError = func(code int, err error) {
		if code == 500 {
			notifyAboutError(code, err)
		}
	}
	// when response sending error.
	var onSendError = func(code int, err error) {
		h.logger.Error("[provide/http]: failed to send response. Error:" + err.Error())
		notifyAboutError(code, err)
	}
	// create iHTTP.
	_http = iHTTP.New(req, res, h.cookie)
	_http.OnHTTPError(onHTTPError)
	_http.OnSendError(onSendError)
	return _http
}

func (h *HTTP) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var _http = h.getInstance(request, response)
		var ctx = context.WithValue(request.Context(), ctxHTTP, _http)
		*request = *request.WithContext(ctx)
		next.ServeHTTP(response, request)
	})
}