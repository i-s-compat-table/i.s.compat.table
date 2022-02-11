package utils

import (
	"strings"
)

func NormalizeString(input string) (normalized string) {
	normalized = strings.Trim(input, " \t\r\n")
	normalized = strings.ReplaceAll(normalized, `“`, `"`)
	normalized = strings.ReplaceAll(normalized, `”`, `"`)
	normalized = strings.ReplaceAll(normalized, `‘`, `'`)
	normalized = strings.ReplaceAll(normalized, `’`, `'`)
	normalized = strings.ReplaceAll(normalized, "\u00A0", ` `) // nonbreaking space
	return normalized
}
