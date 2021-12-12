package core

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// Utils - useful utilities.
type Utils struct {
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
