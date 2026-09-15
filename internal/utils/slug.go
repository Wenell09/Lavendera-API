package utils

import (
	"regexp"
	"strings"
)

var nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)

func GenerateSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = nonAlphaNumeric.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	return value
}
