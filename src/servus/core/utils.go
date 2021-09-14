package core

import (
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

// RemoveSpaces remove spaces from string
func (u *BasicUtils) RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// GetExecuteDir get execution directory
func (u *BasicUtils) GetExecuteDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

// SetCookie set cookie
func (u *BasicUtils) SetCookie(response *http.ResponseWriter, name string, value string) {
	var maxAge, err = u.ConvertTimeWord(Config.Security.Cookie.MaxAge)
	if err != nil {
		Logger.Panic(errors.New("Cookie wrong time. Check your config file."))
	}
	maxAgeSeconds := int(maxAge.Seconds())
	var domain = Config.Security.Cookie.Domain
	var path = Config.Security.Cookie.Path
	var httpOnly = Config.Security.Cookie.HttpOnly
	var secure = Config.Security.Cookie.Secure
	sameSite, err := convertCookieSameSite(Config.Security.Cookie.SameSite)
	if err != nil {
		Logger.Panic(err)
	}
	var cookie = http.Cookie{Name: name, Value: value, Path: path, Domain: domain, MaxAge: maxAgeSeconds, HttpOnly: httpOnly, Secure: secure, SameSite: sameSite}
	http.SetCookie(*response, &cookie)
}

// ConvertTimeWord convert time like "2h"; "2min"; "2sec" to duration
func (u *BasicUtils) ConvertTimeWord(timeShortcut string) (time.Duration, error) {
	timeShortcut = strings.ToLower(timeShortcut)
	timeDuration, err := time.ParseDuration(timeShortcut)
	if err != nil {
		Logger.Panic(errors.New("time converting failed. Is string with time correct?"))
	}
	return timeDuration, nil
}
