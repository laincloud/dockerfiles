package util

import (
	"strings"
)

// EscapeColon replace `:` with `_`
func EscapeColon(s string) string {
	return strings.Replace(s, ":", "_", -1)
}
