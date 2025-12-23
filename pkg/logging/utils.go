package logging

import (
	"strconv"
	"strings"
	"unicode"
)

// EscapeString escapes control characters and non-printable characters in a string
// for safe logging. This prevents user input from breaking log formatting or containing
// malicious control sequences.
func EscapeString(s string) string {
	if s == "" {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(s) * 2)

	for _, r := range s {
		switch r {
		case '\n':
			builder.WriteString("\\n")
		case '\r':
			builder.WriteString("\\r")
		case '\t':
			builder.WriteString("\\t")
		case '\b':
			builder.WriteString("\\b")
		case '\f':
			builder.WriteString("\\f")
		case '\v':
			builder.WriteString("\\v")
		case '\\':
			builder.WriteString("\\\\")
		case '\x00':
			builder.WriteString("\\x00")
		default:
			if unicode.IsPrint(r) {
				builder.WriteRune(r)
			} else if r < 256 {
				builder.WriteString("\\x")
				builder.WriteString(strings.ToUpper(strconv.FormatInt(int64(r), 16)))
			} else {
				builder.WriteString("\\u")
				builder.WriteString(strings.ToUpper(strconv.FormatInt(int64(r), 16)))
			}
		}
	}

	return builder.String()
}

const (
	DefaultTruncateLength = 50
)

// TruncateString limits the length of a string to maxLength runes.
// If the string is longer than maxLength, it is truncated and "..." is appended.
// If maxLength is less than 4, the string is truncated without "..." suffix.
func TruncateString(s string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}

	if maxLength < 4 {
		return string(runes[:maxLength])
	}

	return string(runes[:maxLength-3]) + "..."
}

// Sanitize escapes control characters and truncates a string for safe logging.
// It first escapes the string using EscapeString, then truncates it to DefaultTruncateLength.
// This is a convenience function for the most common logging sanitization needs.
func Sanitize(s string) string {
	return SanitizeWithLength(s, DefaultTruncateLength)
}

// SanitizeWithLength escapes control characters and truncates a string to the specified length for safe logging.
// It first escapes the string using EscapeString, then truncates it to maxLength using TruncateString.
// The escaping is applied first since it may increase the string length.
func SanitizeWithLength(s string, maxLength int) string {
	escaped := EscapeString(s)
	return TruncateString(escaped, maxLength)
}
