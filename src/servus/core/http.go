package core

import (
	"github.com/pkg/errors"
	"net/http"
	"os"
)

type HTTP struct {
	config *ConfigFile
	utils  *Utils
	Logger Logger
}

// Send - wrapper for http.ResponseWriter. Sends response and clear errorMan errors.
func (h *HTTP) Send(response http.ResponseWriter, body string, statusCode int) {
	response.WriteHeader(statusCode)
	_, err := response.Write([]byte(body))
	if err != nil {
		h.Logger.Error("HTTP: failed to send response. Error:" + err.Error())
	}
}

// SetCookie - cool wrapper for setting cookies.
func (h *HTTP) SetCookie(response *http.ResponseWriter, name string, value string) {
	var maxAge, err = h.utils.ConvertTimeWord(h.config.Security.Cookie.MaxAge)
	if err != nil {
		var errPretty = errors.Wrap(err, "core: SetCookie convert time failed. Error")
		h.Logger.Panic(errPretty)
		os.Exit(1)
	}
	maxAgeSeconds := int(maxAge.Seconds())
	var domain = h.config.Security.Cookie.Domain
	var path = h.config.Security.Cookie.Path
	var httpOnly = h.config.Security.Cookie.HttpOnly
	var secure = h.config.Security.Cookie.Secure
	sameSite, err := convertCookieSameSite(h.config.Security.Cookie.SameSite)
	if err != nil {
		var errPretty = errors.Wrap(err, "core: SetCookie failed convert cookie sameSite. Error")
		h.Logger.Panic(errPretty)
		os.Exit(1)
	}
	var cookie = &http.Cookie{Name: name, Value: value, Path: path, Domain: domain, MaxAge: maxAgeSeconds, HttpOnly: httpOnly, Secure: secure, SameSite: sameSite}
	http.SetCookie(*response, cookie)
}
