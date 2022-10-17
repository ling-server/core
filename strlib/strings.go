package strlib

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TrimsLineBreaks trims line breaks in string.
func TrimLineBreaks(s string) string {
	escaped := strings.ReplaceAll(s, "\n", "")
	escaped = strings.ReplaceAll(escaped, "\r", "")
	return escaped
}

// Title uppercase the first character, and lower case the rest, for example covert MANUAL to Manual
func Title(s string) string {
	title := cases.Title(language.Und)
	return title.String(strings.ToLower(s))
}
