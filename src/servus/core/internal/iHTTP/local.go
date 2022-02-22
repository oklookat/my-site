package iHTTP

import (
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (i *Instance) wrapError(err error) error {
	if err == nil {
		return nil
	}
	err = errors.Wrap(err, "[iHTTP] ")
	return err
}

func (i *Instance) newError(message string) error {
	return errors.New("[iHTTP] " + message)
}

// convert cookie sameSite string to http.SameSite.
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
		return http.SameSiteDefaultMode, i.newError("wrong sameSite string")
	}
}

// convert time like "2h"; "2min"; "2sec" to duration.
func (i *Instance) convertTimeWord(timeShortcut string) (dur time.Duration, err error) {
	timeShortcut = strings.ToLower(timeShortcut)
	dur, err = time.ParseDuration(timeShortcut)
	if err != nil {
		err = i.wrapError(err)
		return
	}
	return
}
