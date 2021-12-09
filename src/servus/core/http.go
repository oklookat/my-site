package core

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"servus/core/internal/zipify"
	"strings"
)

// HTTP - cool things for request/response.
type HTTP struct {
	core     *Core
	request  *http.Request
	response http.ResponseWriter
}

// ProvideHTTP - gives access to core.HTTP in request context.
func (m *Middleware) ProvideHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var h = &HTTP{core: m.core, request: request, response: response}
		var ctx = context.WithValue(request.Context(), ctxHTTP, h)
		*request = *request.WithContext(ctx)
		next.ServeHTTP(response, request)
	})
}

// GetHTTP - get HTTP from request context.
func (m *Middleware) GetHTTP(request *http.Request) *HTTP {
	var ctx = request.Context()
	var h, ok = ctx.Value(ctxHTTP).(*HTTP)
	if !ok {
		m.core.Logger.Warn("middleware: error while get HTTP struct from request context. Are you provided ProvideHTTP middleware?")
		return nil
	}
	return h
}

func (h *HTTP) new(core *Core, req *http.Request, res http.ResponseWriter) {
	h.core = core
	h.request = req
	h.response = res
}

// Send - sends response and log it if error.
func (h *HTTP) Send(body string, statusCode int, err error) {
	h.response.WriteHeader(statusCode)


	// TODO: add 500 error handler here; notification about 500 error with useful information
	// TODO: get stack, request path, write in txt and zip, send to telegram
	if statusCode == 500 {
		var stack = h.core.Utils.GetStackTrace(err)
		stack.trace = "ERROR 500. TRACE:\n" + stack.trace
		var arc = zipify.ZipFile{}
		err = arc.AddFile("dump-" + stack.timestamp, &stack)
		if err == nil {
			var caption = "500 error"
			h.core.Control.tg.sendFile(&caption, stack.timestamp, arc.GetFile())
		}
	}


	_, err = h.response.Write([]byte(body))
	if err != nil {
		h.core.Logger.Error("HTTP: failed to send response. Error:" + err.Error())
	}
}

// SetCookie - set cookie.
func (h *HTTP) SetCookie(name string, value string) {
	var cfg = h.core.Config
	var maxAge, err = h.core.Utils.ConvertTimeWord(cfg.Security.Cookie.MaxAge)
	if err != nil {
		var errPretty = errors.Wrap(err, "HTTP: SetCookie convert time failed. Error")
		h.core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	maxAgeSeconds := int(maxAge.Seconds())
	var domain = cfg.Security.Cookie.Domain
	var path = cfg.Security.Cookie.Path
	var httpOnly = cfg.Security.Cookie.HttpOnly
	var secure = cfg.Security.Cookie.Secure
	sameSite, err := h.convertCookieSameSite(cfg.Security.Cookie.SameSite)
	if err != nil {
		var errPretty = errors.Wrap(err, "HTTP: SetCookie failed convert cookie sameSite. Error")
		h.core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	var cookie = &http.Cookie{Name: name, Value: value, Path: path, Domain: domain, MaxAge: maxAgeSeconds, HttpOnly: httpOnly, Secure: secure, SameSite: sameSite}
	http.SetCookie(h.response, cookie)
}

// convertCookieSameSite - convert cookie sameSite string to http.SameSite.
func (h *HTTP) convertCookieSameSite(sameSite string) (http.SameSite, error) {
	sameSite = strings.ToUpper(sameSite)
	switch sameSite {
	case "DEFAULT":
		return http.SameSiteDefaultMode, nil
	case "LAX":
		return http.SameSiteLaxMode, nil
	case "STRICT":
		return http.SameSiteStrictMode, nil
	case "NONE":
		return http.SameSiteNoneMode, nil
	default:
		return http.SameSiteDefaultMode, errors.New("HTTP: wrong sameSite string")
	}
}
