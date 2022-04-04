package validator

import (
	"net"
	"regexp"
	"strings"
)

// string with spaces only anyway will be empty.
func IsEmpty(text *string) bool {
	if text == nil {
		return true
	}
	return len(strings.TrimSpace(*text)) == 0
}

// a-z A-Z 0-9
func IsAlphanumeric(text *string) bool {
	if text == nil {
		return false
	}
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	return alphanumeric.MatchString(*text)
}

// alphanumeric and keyboard symbols?
func IsAlphanumericWithSymbols(text *string) bool {
	if text == nil {
		return false
	}
	alphaWithSymbols := regexp.MustCompile(`^[a-zA-Z0-9\-+~"'\\x60(){\\[}|:;,.!=@#$%^&?№*\\\\\\\\/<>]*$`)
	return alphaWithSymbols.MatchString(*text)
}

// IPv4 or IPv6 address?
func IsIP(text *string) bool {
	if text == nil {
		return false
	}
	var parsed = net.ParseIP(*text)
	return parsed != nil
}

// string violates min and max length constraints?
func MinMax(text *string, min int, max int) bool {
	if text == nil {
		return false
	}
	var textLen = len([]rune(*text))
	return textLen < min || textLen > max
}
