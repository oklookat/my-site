package core

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/oklog/ulid/v2"
)

type utils struct {
}

// remove spaces from string.
func (u *utils) RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// fromat path slashes.
func (u *utils) FormatPath(path string) string {
	path = filepath.ToSlash(path)
	path = filepath.Clean(path)
	return path
}

// get HTTP from request context.
func (u *utils) GetHTTP(request *http.Request) HTTP {
	var ctx = request.Context()
	var h, _ = ctx.Value(ctxHTTP).(HTTP)
	return h
}

// generate ULID.
func (u *utils) GenerateULID() (ul string, err error) {
	current := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(current.UnixNano())), 0)
	ulType, err := ulid.New(ulid.Timestamp(current), entropy)
	if err != nil {
		return "", err
	}
	ul = ulType.String()
	return
}

// get rune length.
func (u *utils) LenRune(val string) int {
	return len([]rune(val))
}

// callback can only be called once every 5 seconds.
func (u *utils) Debounce(interval time.Duration) (debouncer func(callback func())) {
	var isCooldown = false
	return func(callback func()) {
		if isCooldown {
			return
		}
		if callback != nil {
			callback()
		}
		isCooldown = true
		time.AfterFunc(interval, func() { isCooldown = false })
	}
}
