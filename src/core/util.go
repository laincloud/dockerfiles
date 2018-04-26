package core

import (
	"strings"
)

// EscapeSlash replace `/` with `_`
func EscapeSlash(s string) string {
	return strings.Replace(s, "/", "_", -1)
}
