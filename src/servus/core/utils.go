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

func (u *utils) RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (u *utils) GetExecutionDir() (path string, err error) {
	path, err = os.Executable()
	if err != nil {
		err = errors.Wrap(err, "[core/utils]: failed to get execution directory. Error")
		return
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		err = errors.Wrap(err, "[core/utils]: failed to get execution directory / symlink follow. Error")
		return
	}
	return
}

func (u *utils) FormatPath(path string) string {
	path = filepath.FromSlash(path)
	path = filepath.ToSlash(path)
	return path
}

func (u *utils) GetHTTP(request *http.Request) HTTP {
	var ctx = request.Context()
	var h, _ = ctx.Value(ctxHTTP).(HTTP)
	return h
}

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
