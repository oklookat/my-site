package core

import (
	"github.com/pkg/errors"
	"net/http"
	"os"
	"servus/core/internal/corsParse"
	"strings"
)

type HTTP struct {
	config     *ConfigFile
	utils      *Utils
	logger     Logger
	Middleware *Middleware
}


// bootHTTP - boot http utilities. Use this after booting the config.
func (c *Core) bootHTTP() {
	var corsConfig = corsParse.Config{
		AllowCredentials: c.Config.Security.CORS.AllowCredentials,
		AllowOrigin:      c.Config.Security.CORS.AllowOrigin,
		AllowMethods:     c.Config.Security.CORS.AllowMethods,
		AllowHeaders:     c.Config.Security.CORS.AllowHeaders,
		ExposeHeaders:    c.Config.Security.CORS.ExposeHeaders,
		MaxAge:           c.Config.Security.CORS.MaxAge,
	}
	var corsInstance = corsParse.New(corsConfig)
	middleware := Middleware{config: c.Config, cors: &corsInstance}
	c.HTTP = &HTTP{config: c.Config, Middleware: &middleware}
}

// Send - wrapper for http.ResponseWriter. Sends response and clear errorMan errors.
func (h *HTTP) Send(response http.ResponseWriter, body string, statusCode int) {
	response.WriteHeader(statusCode)
	// TODO: add 500 error handler here; notification about 500 error with useful information
	// TODO: get stack, request path, write in txt and zip, send to telegram
	_, err := response.Write([]byte(body))
	if err != nil {
		h.logger.Error("HTTP: failed to send response. Error:" + err.Error())
	}
}

// SetCookie - cool wrapper for setting cookies.
func (h *HTTP) SetCookie(response *http.ResponseWriter, name string, value string) {
	var maxAge, err = h.utils.ConvertTimeWord(h.config.Security.Cookie.MaxAge)
	if err != nil {
		var errPretty = errors.Wrap(err, "HTTP: SetCookie convert time failed. Error")
		h.logger.Panic(errPretty)
		os.Exit(1)
	}
	maxAgeSeconds := int(maxAge.Seconds())
	var domain = h.config.Security.Cookie.Domain
	var path = h.config.Security.Cookie.Path
	var httpOnly = h.config.Security.Cookie.HttpOnly
	var secure = h.config.Security.Cookie.Secure
	sameSite, err := h.convertCookieSameSite(h.config.Security.Cookie.SameSite)
	if err != nil {
		var errPretty = errors.Wrap(err, "HTTP: SetCookie failed convert cookie sameSite. Error")
		h.logger.Panic(errPretty)
		os.Exit(1)
	}
	var cookie = &http.Cookie{Name: name, Value: value, Path: path, Domain: domain, MaxAge: maxAgeSeconds, HttpOnly: httpOnly, Secure: secure, SameSite: sameSite}
	http.SetCookie(*response, cookie)
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
