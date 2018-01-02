package util

import (
	"testing"
)

func TestEscapeColon(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"laincloud/centos:7", "laincloud/centos_7"},
	}

	for _, c := range cases {
		got := EscapeColon(c.in)
		if got != c.want {
			t.Errorf("EscapeColon(%s) == %s, want %s.", c.in, got, c.want)
		}
	}
}
