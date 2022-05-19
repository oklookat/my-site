package iHTTP

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"errors"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	err = fmt.Errorf("[iHTTP] %w", err)
	return err
}

func newError(message string) error {
	return errors.New("[iHTTP] " + message)
}

// convert cookie sameSite string to http.SameSite.
func convertCookieSameSite(sameSite string) (http.SameSite, error) {
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
		return http.SameSiteDefaultMode, newError("wrong sameSite string")
	}
}

// convert time like "2h"; "2min"; "2sec" to duration.
func convertTimeWord(timeShortcut string) (dur time.Duration, err error) {
	timeShortcut = strings.ToLower(timeShortcut)
	dur, err = time.ParseDuration(timeShortcut)
	if err != nil {
		err = wrapError(err)
		return
	}
	return
}
