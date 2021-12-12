package iHTTP

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// Instance - cool things for request/response.
type Instance struct {
	cookie      *ConfigCookie
	onHTTPError func(code int, err error)
	onSendError func(code int, err error)
	request     *http.Request
	response    http.ResponseWriter
}

// ConfigCookie - config for setting cookies.
type ConfigCookie struct {
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	MaxAge   string `json:"maxAge"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"sameSite"`
}

// New - creates new HTTP instance.
func New(req *http.Request, res http.ResponseWriter, cookie *ConfigCookie) *Instance {
	var i = &Instance{}
	i.request = req
	i.response = res
	i.cookie = cookie
	return i
}

// OnHTTPError - when 299+ error.
func (i *Instance) OnHTTPError(callback func(code int, err error)) {
	i.onHTTPError = callback
}

// OnSendError - response sending error.
func (i *Instance) OnSendError(callback func(code int, err error)) {
	i.onSendError = callback
}

// Send - sends response and log it if error.
func (i *Instance) Send(body string, statusCode int, err error) {
	if err != nil {
		err = errors.Wrap(err, "")
	}
	var is1xx = (statusCode-100) >= 0 && (statusCode-100) <= 99
	var is2xx = (statusCode-200) >= 0 && (statusCode-200) <= 99
	var is3xx = (statusCode-300) >= 0 && (statusCode-300) <= 99
	// http error code and callback not empty?
	if !is1xx && !is2xx && !is3xx && i.onHTTPError != nil {
		go i.onHTTPError(statusCode, err)
	}
	i.response.WriteHeader(statusCode)
	_, err = i.response.Write([]byte(body))
	if err != nil && i.onSendError != nil {
		go i.onSendError(statusCode, err)
	}
}

// SetCookie - set cookie.
func (i *Instance) SetCookie(name string, value string) error {
	var maxAge, err = i.convertTimeWord(i.cookie.MaxAge)
	if err != nil {
		var errPretty = errors.Wrap(err, "[http/SetCookie]: convert time failed. Error: ")
		return errPretty
	}
	maxAgeSeconds := int(maxAge.Seconds())
	var domain = i.cookie.Domain
	var path = i.cookie.Path
	var httpOnly = i.cookie.HttpOnly
	var secure = i.cookie.Secure
	sameSite, err := i.convertCookieSameSite(i.cookie.SameSite)
	if err != nil {
		var errPretty = errors.Wrap(err, "[http/SetCookie]: failed convert sameSite. Error: ")
		return errPretty
	}
	var cookie = &http.Cookie{Name: name, Value: value, Path: path, Domain: domain, MaxAge: maxAgeSeconds, HttpOnly: httpOnly, Secure: secure, SameSite: sameSite}
	http.SetCookie(i.response, cookie)
	return err
}

// convertCookieSameSite - convert cookie sameSite string to http.SameSite.
func (i *Instance) convertCookieSameSite(sameSite string) (http.SameSite, error) {
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

// GetDump - get pretty HTTP request info in string/io.Reader.
func (i *Instance) GetDump() io.Reader {
	// cookies.
	var cookies = i.request.Cookies()
	var cookieString = ""
	for _, cookie := range cookies {
		var cook = fmt.Sprintf(`
////////////
domain: %v
path: %v
name: %v
value: %v
secure: %v
////////////

        `, cookie.Domain, cookie.Path, cookie.Name, cookie.Value, cookie.Secure)
		cookieString = cookieString + cook
	}
	if len(cookieString) < 2 {
		cookieString = "No cookies.\n"
	}
	// url.
	var url = fmt.Sprintf(`
method: %v
path: %v
rawPath: %v
contentLength: %v
rawQuery: %v
	`, i.request.Method, i.request.URL.Path, i.request.URL.RawPath, i.request.ContentLength, i.request.URL.RawQuery)
	// total.
	var dump = fmt.Sprintf(`----URL----%v

----Cookies----%v`, url, cookieString)
	return strings.NewReader(dump)
}

// convertTimeWord - convert time like "2h"; "2min"; "2sec" to duration.
func (i *Instance) convertTimeWord(timeShortcut string) (time.Duration, error) {
	timeShortcut = strings.ToLower(timeShortcut)
	timeDuration, err := time.ParseDuration(timeShortcut)
	if err != nil {
		var errPretty = errors.Wrap(err, "[http/convertTimeWord]: converting failed. Error: ")
		return 0, errPretty
	}
	return timeDuration, nil
}
