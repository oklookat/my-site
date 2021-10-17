package core

import (
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
)

type HTTP struct {
	config     *ConfigFile
	utils      *Utils
	logger     Logger
	Middleware *Middleware
}

// Send - wrapper for http.ResponseWriter. Sends response and clear errorMan errors.
func (h *HTTP) Send(response http.ResponseWriter, body string, statusCode int) {
	response.WriteHeader(statusCode)
	_, err := response.Write([]byte(body))
	if err != nil {
		h.logger.Error("HTTP: failed to send response. Error:" + err.Error())
	}
}

// SetCookie - cool wrapper for setting cookies.
func (h *HTTP) SetCookie(response *http.ResponseWriter, name string, value string) {
	var maxAge, err = h.utils.ConvertTimeWord(h.config.Security.Cookie.MaxAge)
	if err != nil {
		var errPretty = errors.Wrap(err, "core: SetCookie convert time failed. Error")
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
		var errPretty = errors.Wrap(err, "core: SetCookie failed convert cookie sameSite. Error")
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
		return http.SameSiteDefaultMode, errors.New("Wrong sameSite string.")
	}
}

// Middleware - basic middleware.
type Middleware struct {
	config *ConfigFile
	cors   cors
}

// AsJSON - this middleware set application/json header.
func (m *Middleware) AsJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(writer, request)
	})
}

// CORS - CORS middleware depending on config file.
func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if m.config.Security.CORS.Active {
			var isPreflight = m.cors.SetHeaders(writer, request)
			if isPreflight {
				return
			}
		}
		next.ServeHTTP(writer, request)
	})
}

// Security - basic security.
func (m *Middleware) Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var method = request.Method
		method = strings.ToUpper(method)
		//////// request body size check
		if m.config.Security.Limiter.Body.Active {
			var isNotReadOnlyMethod = !(method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions)
			if isNotReadOnlyMethod {
				var requestURI = request.RequestURI
				var exceptSlice = m.config.Security.Limiter.Body.Except
				var isBypassed = false
				for _, exceptOne := range exceptSlice {
					if exceptOne == requestURI {
						isBypassed = true
					}
				}
				if !isBypassed {
					request.Body = http.MaxBytesReader(writer, request.Body, m.config.Security.Limiter.Body.MaxSize)
					//payload, err := ioutil.ReadAll(request.Body)
					//if err != nil {
					//	if strings.Contains(err.Error(), "request body too large"){
					//		http.Error(writer, "Request body too large.", 413)
					//		return
					//	}
					//}
				}
			}

		}
		next.ServeHTTP(writer, request)
	})
}
