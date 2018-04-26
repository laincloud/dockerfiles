package core

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
		{" FROM laincloud/centos:7", FromImage{
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

func TestParseTags(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"# TAGS 7.4.1708 7 latest", []string{
			"7.4.1708",
			"7",
			"latest",
		}},
		{"#TAGS 7.4.1708 7 latest", []string{
			"7.4.1708",
			"7",
			"latest",
		}},
		{" # TAGS 7.4.1708 7 latest", []string{
			"7.4.1708",
			"7",
			"latest",
		}},
		{"# TAGS 7.4.1708  7  latest", []string{
			"7.4.1708",
			"7",
			"latest",
		}},
	}

	for _, c := range cases {
		got, err := parseTags(c.in)
		if err != nil {
			t.Errorf("parseTags(%s) failed, error: %s.", c.in, err)
		}

		for i, tag := range got {
			if tag != c.want[i] {
				t.Errorf("parseTags(%s) == %+v, want %+v.", c.in, got, c.want)
			}
		}
	}
}
