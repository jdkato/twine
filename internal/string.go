package internal

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// StringInSlice determines if `slice` contains the string `a`.
func StringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if a == b {
			return true
		}
	}
	return false
}

// HasAnySuffix determines if `s` has any of the suffixes in `suffixes`.
func HasAnySuffix(s string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func HasAnyPrefix(s string, prefixes []string) bool {
	n := len(s)
	for _, prefix := range prefixes {
		if n > len(prefix) && strings.HasPrefix(s, prefix) {
			return true
		}

	}
	return false
}

func HasAnyIndex(s string, suffixes []string) int {
	n := len(s)
	for _, suffix := range suffixes {
		idx := strings.Index(s, suffix)
		if idx >= 0 && n > len(suffix) {
			return idx
		}
	}
	return -1
}

// CharAt returns the ith character of s, if it exists. Otherwise, it returns
// the first character.
func CharAt(s string, i int) byte {
	if i >= 0 && i < len(s) {
		return s[i]
	}
	return s[0]
}

// ToTitle returns a copy of the string m with its first Unicode letter mapped
// to its title case.
func ToTitle(m string) string {
	r, size := utf8.DecodeRuneInString(m)
	return string(unicode.ToTitle(r)) + strings.ToLower(m[size:])
}
