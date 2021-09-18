package validator

import (
	"regexp"
	"strings"
)

// IsEmpty string with spaces only anyway will be empty
func IsEmpty(text string) bool{
	return len(strings.TrimSpace(text)) == 0
}

// IsAlphanumeric a-z A-Z 0-9
func IsAlphanumeric(text string) bool {
	alphanumeric := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return alphanumeric.MatchString(text)
}

// IsAlphanumericWithSymbols a-z A-Z 0-9 _ @ ! . / # & + * %
func IsAlphanumericWithSymbols(text string) bool {
	alphaWithSymbols := regexp.MustCompile("^[a-zA-Z0-9_@!./#&+*%-]*$")
	return alphaWithSymbols.MatchString(text)
}
