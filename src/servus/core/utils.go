package core

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// Utils - must have utilities.
type Utils struct {
	config *ConfigFile
	logger Logger
}

// RemoveSpaces - remove spaces from string.
func (u *Utils) RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// GetExecutionDir - get server execution directory.
func (u *Utils) GetExecutionDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.ToSlash(path)
}

// FormatPath - format path to system specific slashes.
func (u *Utils) FormatPath(path string) string {
	path = filepath.FromSlash(path)
	path = filepath.ToSlash(path)
	return path
}

// ConvertTimeWord - convert time like "2h"; "2min"; "2sec" to duration (uses time.ParseDuration).
func (u *Utils) ConvertTimeWord(timeShortcut string) (time.Duration, error) {
	timeShortcut = strings.ToLower(timeShortcut)
	timeDuration, err := time.ParseDuration(timeShortcut)
	if err != nil {
		var errPretty = errors.Wrap(err, "core: convertTimeWord time converting failed. Error")
		u.logger.Panic(errPretty)
		os.Exit(1)
	}
	return timeDuration, nil
}
