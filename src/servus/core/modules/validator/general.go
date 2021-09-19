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

// IsAlphanumericWithSymbols alphanumeric and keyboard symbols
func IsAlphanumericWithSymbols(text string) bool {
	alphaWithSymbols := regexp.MustCompile("^[a-zA-Z0-9~`!@#$%^&*()_\\-+={\\[}\\]|:;\"'<,>.?\\/\\\\\\\\]*$")
	return alphaWithSymbols.MatchString(text)
}
