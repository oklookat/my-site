package core

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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
