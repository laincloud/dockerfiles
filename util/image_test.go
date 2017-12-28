package util

import (
	"testing"
)

func TestParseFromImage(t *testing.T) {
	cases := []struct {
		in   string
		want FromImage
	}{
		{"FROM laincloud/centos:7", FromImage{
			Organization: "laincloud",
			Repository:   "centos",
			Tag:          "7",
		}},
		{"From  laincloud/centos:7", FromImage{
			Organization: "laincloud",
			Repository:   "centos",
			Tag:          "7",
		}},
		{"FROM centos", FromImage{
			Organization: "library",
			Repository:   "centos",
			Tag:          "latest",
		}},
	}

	for _, c := range cases {
		got, err := parseFromImage(c.in)
		if err != nil {
			t.Errorf("parseFromImage(%s) failed, error: %s.", c.in, err)
		}
		if *got != c.want {
			t.Errorf("parseFromImage(%s) == %s, want %s.", c.in, got, c.want)
		}
	}
}
