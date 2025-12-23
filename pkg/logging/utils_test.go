package logging

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "regular text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "newline character",
			input:    "Hello\nWorld",
			expected: "Hello\\nWorld",
		},
		{
			name:     "carriage return",
			input:    "Hello\rWorld",
			expected: "Hello\\rWorld",
		},
		{
			name:     "tab character",
			input:    "Hello\tWorld",
			expected: "Hello\\tWorld",
		},
		{
			name:     "backspace character",
			input:    "Hello\bWorld",
			expected: "Hello\\bWorld",
		},
		{
			name:     "form feed character",
			input:    "Hello\fWorld",
			expected: "Hello\\fWorld",
		},
		{
			name:     "vertical tab character",
			input:    "Hello\vWorld",
			expected: "Hello\\vWorld",
		},
		{
			name:     "backslash character",
			input:    "Hello\\World",
			expected: "Hello\\\\World",
		},
		{
			name:     "null character",
			input:    "Hello\x00World",
			expected: "Hello\\x00World",
		},
		{
			name:     "multiple control characters",
			input:    "Hello\n\r\t\bWorld",
			expected: "Hello\\n\\r\\t\\bWorld",
		},
		{
			name:     "unicode characters",
			input:    "Hello 世界",
			expected: "Hello 世界",
		},
		{
			name:     "control character in hex range",
			input:    "Hello\x01World",
			expected: "Hello\\x1World",
		},
		{
			name:     "bell character",
			input:    "Hello\x07World",
			expected: "Hello\\x7World",
		},
		{
			name:     "delete character",
			input:    "Hello\x7fWorld",
			expected: "Hello\\x7FWorld",
		},
		{
			name:     "emoji characters",
			input:    "Hello 👋 World",
			expected: "Hello 👋 World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "empty string",
			input:     "",
			maxLength: 10,
			expected:  "",
		},
		{
			name:      "string shorter than max",
			input:     "Hello",
			maxLength: 10,
			expected:  "Hello",
		},
		{
			name:      "string equal to max",
			input:     "Hello",
			maxLength: 5,
			expected:  "Hello",
		},
		{
			name:      "string longer than max",
			input:     "Hello World",
			maxLength: 8,
			expected:  "Hello...",
		},
		{
			name:      "very long string",
			input:     "This is a very long string that should be truncated",
			maxLength: 20,
			expected:  "This is a very lo...",
		},
		{
			name:      "max length zero",
			input:     "Hello",
			maxLength: 0,
			expected:  "",
		},
		{
			name:      "negative max length",
			input:     "Hello",
			maxLength: -5,
			expected:  "",
		},
		{
			name:      "max length 1",
			input:     "Hello",
			maxLength: 1,
			expected:  "H",
		},
		{
			name:      "max length 2",
			input:     "Hello",
			maxLength: 2,
			expected:  "He",
		},
		{
			name:      "max length 3",
			input:     "Hello",
			maxLength: 3,
			expected:  "Hel",
		},
		{
			name:      "max length 4 exactly",
			input:     "Hello",
			maxLength: 4,
			expected:  "H...",
		},
		{
			name:      "unicode string truncation",
			input:     "Hello 世界 World",
			maxLength: 10,
			expected:  "Hello 世...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateString(tt.input, tt.maxLength)
			assert.Equal(t, tt.expected, result)
			assert.True(t, len([]rune(result)) <= tt.maxLength || tt.maxLength <= 0)
		})
	}
}

func TestEscapeStringPerformance(t *testing.T) {
	longString := strings.Repeat("Hello World\n\r\t", 1000)
	result := EscapeString(longString)
	assert.Contains(t, result, "\\n")
	assert.Contains(t, result, "\\r")
	assert.Contains(t, result, "\\t")
}

func TestTruncateStringPerformance(t *testing.T) {
	longString := strings.Repeat("a", 10000)
	result := TruncateString(longString, 100)
	assert.Equal(t, 100, len([]rune(result)))
	assert.True(t, strings.HasSuffix(result, "..."))
}

func TestCombinedFunctions(t *testing.T) {
	input := "User message with\nnewlines and\ttabs that is quite long and should be both escaped and truncated"
	escaped := EscapeString(input)
	truncated := TruncateString(escaped, 50)

	assert.Contains(t, truncated, "\\n")
	assert.Contains(t, truncated, "\\t")
	assert.True(t, len([]rune(truncated)) <= 50)
	assert.True(t, strings.HasSuffix(truncated, "..."))
}
