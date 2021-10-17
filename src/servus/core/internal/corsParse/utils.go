package corsParse

import (
	"strings"
	"unicode"
)

// removeSpaces - remove spaces from string.
func removeSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}