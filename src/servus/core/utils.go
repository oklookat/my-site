package core

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
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

// get dir where bin executes.
func (u *utils) GetExecutionDir() (path string, err error) {
	path, err = os.Executable()
	if err != nil {
		err = errors.Wrap(err, "[core/utils]: failed to get execution directory. Error")
		return
	}
	path, err = filepath.Abs(path)
	if err != nil {
		err = errors.Wrap(err, "[core/utils]: failed to get execution directory / absolute path. Error")
		return
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		err = errors.Wrap(err, "[core/utils]: failed to get execution directory / symlink follow. Error")
		return
	}
	return
}

// fromat path slashes.
func (u *utils) FormatPath(path string) string {
	path = filepath.FromSlash(path)
	path = filepath.ToSlash(path)
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
