package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"servus/core/internal/iHTTP"
	"servus/core/internal/zipify"
	"time"

	"github.com/oklookat/goway"
)

type _ctxHTTP string

const ctxHTTP _ctxHTTP = "SERVUS_HTTP_HELPER"

func (i *Instance) setupHTTP() error {
	var helper = httpHelper{}
	var err = helper.new(i.Logger, i.Control, i.Config.Security.Cookie)
	if err != nil {
		return err
	}
	i.Http = helper
	return nil
}

type httpHelper struct {
	logger  Logger
	control Controller
	cookie  *iHTTP.ConfigCookie
}

func (h *httpHelper) new(
	log Logger,
	control Controller,
	cookie *iHTTP.ConfigCookie) error {

	if log == nil {
		return errors.New("[HTTP] logger nil pointer")
	}
	if control == nil {
		return errors.New("[HTTP] controller nil pointer")
	}
	if cookie == nil {
		return errors.New("[HTTP] cookie nil pointer")
	}

	h.logger = log
	h.control = control
	h.cookie = cookie

	return nil
}

func (h *httpHelper) getInstance(req *http.Request, res http.ResponseWriter) *iHTTP.Instance {
	return &iHTTP.Instance{
		Request:         req,
		Response:        res,
		Cookie:          h.cookie,
		OnHTTPError:     h.onHTTPError,
		OnSendError:     h.onResponseSendError,
		RouteArgsGetter: goway.Vars,
	}
}

func (h *httpHelper) Get(request *http.Request) HTTP {
	var ctx = request.Context()
	var theHttp, _ = ctx.Value(ctxHTTP).(HTTP)
	return theHttp
}

func (h *httpHelper) onHTTPError(theHTTP *iHTTP.Instance, code int, err error) {
	if code != 500 || theHTTP == nil {
		return
	}

	if err == nil {
		err = errors.New("*empty error*")
	}

	// log.
	h.logger.Error("[core/http] " + err.Error())

	// create zip.
	var zip = zipify.New()
	err = zip.AddFile("request.txt", theHTTP.GetDump())

	// if zip error.
	if err != nil {
		var message = fmt.Sprintf("[#ERROR #code%v] make .zip also failed.", code)
		h.logger.Error("[core/http]: " + err.Error())
		h.control.SendMessage(message)
		return
	}

	// send dump.
	var message = fmt.Sprintf("[#ERROR #code%v]", code)
	var timestamp = fmt.Sprint(time.Now().Unix())
	h.control.SendFile(&message, timestamp+".zip", zip.GetBytesAndClose())
}

func (h *httpHelper) onResponseSendError(theHTTP *iHTTP.Instance, code int, err error) {
	h.logger.Error("[core/http] failed to send response. Error:" + err.Error())
	h.onHTTPError(theHTTP, 500, err)
}

func (h *httpHelper) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var _http = h.getInstance(request, response)
		var ctx = context.WithValue(request.Context(), ctxHTTP, _http)
		*request = *request.WithContext(ctx)
		next.ServeHTTP(response, request)
	})
}
