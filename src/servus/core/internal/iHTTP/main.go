package iHTTP

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// config for setting cookies.
type ConfigCookie struct {
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	MaxAge   string `json:"maxAge"`
	HttpOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
	SameSite string `json:"sameSite"`
}

// Instance - cool things for request/response.
type Instance struct {
	Request *http.Request

	Response http.ResponseWriter

	Cookie *ConfigCookie

	// when 399+ error.
	OnHTTPError func(i *Instance, code int, err error)

	// response sending error.
	OnSendError func(i *Instance, code int, err error)

	// route arguments getter.
	RouteArgsGetter func(r *http.Request) map[string]string
}

// sends response and log if error.
func (i *Instance) Send(body string, statusCode int, err error) {
	if i.Response == nil {
		return
	}

	// is http error and http error callback not empty?
	var isHTTPError = statusCode > 399
	if isHTTPError && i.OnHTTPError != nil {
		go i.OnHTTPError(i, statusCode, err)
	}
	i.Response.WriteHeader(statusCode)
	if _, err = i.Response.Write([]byte(body)); err != nil && i.OnSendError != nil {
		go i.OnSendError(i, statusCode, err)
	}
}

// set cookie.
func (i *Instance) SetCookie(name string, value string) error {
	if i.Cookie == nil || i.Response == nil {
		return nil
	}

	var maxAge, err = convertTimeWord(i.Cookie.MaxAge)
	if err != nil {
		return wrapError(err)
	}

	var maxAgeSeconds = int(maxAge.Seconds())
	var domain = i.Cookie.Domain
	var path = i.Cookie.Path
	var httpOnly = i.Cookie.HttpOnly
	var secure = i.Cookie.Secure

	sameSite, err := convertCookieSameSite(i.Cookie.SameSite)
	if err != nil {
		return wrapError(err)
	}

	var cookie = &http.Cookie{Name: name, Value: value, Path: path, Domain: domain,
		MaxAge: maxAgeSeconds, HttpOnly: httpOnly,
		Secure: secure, SameSite: sameSite}
	http.SetCookie(i.Response, cookie)
	return err
}

func (i *Instance) UnsetCookie(name string) error {
	if i.Cookie == nil || i.Response == nil {
		return nil
	}

	var maxAgeSeconds = 1
	var domain = i.Cookie.Domain
	var path = i.Cookie.Path
	var httpOnly = i.Cookie.HttpOnly
	var secure = i.Cookie.Secure

	sameSite, err := convertCookieSameSite(i.Cookie.SameSite)
	if err != nil {
		return wrapError(err)
	}

	var cookie = &http.Cookie{Name: name, Value: "", Path: path, Domain: domain,
		MaxAge: maxAgeSeconds, HttpOnly: httpOnly,
		Secure: secure, SameSite: sameSite}
	http.SetCookie(i.Response, cookie)
	return err
}

// get pretty HTTP request info in string/io.Reader.
func (i *Instance) GetDump() io.Reader {
	if i.Request == nil {
		return nil
	}

	// cookies.
	var cookies = i.Request.Cookies()
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
`, i.Request.Method, i.Request.URL.Path, i.Request.URL.RawPath,
		i.Request.ContentLength, i.Request.URL.RawQuery)
	// total.
	var dump = fmt.Sprintf(`----URL----%v

----Cookies----%v`, url, cookieString)
	return strings.NewReader(dump)
}

// get route arguments.
//
// like route: /users/{username}/{id}
//
// request: /users/iam/111
//
// and map will be: [username: iam, id: 111]
func (i *Instance) GetRouteArgs() map[string]string {
	if i.Request == nil {
		return nil
	}
	if i.RouteArgsGetter == nil {
		fmt.Println("[iHTTP] warning: GetRouteArgs called, but RouteArgsGetter is nil")
		return nil
	}
	return i.RouteArgsGetter(i.Request)
}
