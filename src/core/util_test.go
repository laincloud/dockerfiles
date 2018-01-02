package util

import (
	"testing"
)

func TestEscapeSlash(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"centos/7", "centos_7"},
	}

	for _, c := range cases {
		got := EscapeSlash(c.in)
		if got != c.want {
			t.Errorf("EscapeColon(%s) == %s, want %s.", c.in, got, c.want)
		}
	}
}
