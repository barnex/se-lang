package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {

	cases := []struct {
		in   string
		want []Token
	}{
		{"", []Token{{TEOF, ""}}},
		{"1", []Token{{TNum, "1"}, {TEOF, ""}}},
		{"23", []Token{{TNum, "23"}, {TEOF, ""}}},
		{" 45\t678 ", []Token{{TNum, "45"}, {TNum, "678"}, {TEOF, ""}}},
	}

	for _, c := range cases {
		have, err := Lex(c.in)
		if err != nil {
			t.Errorf("%q: error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%q: have %v, want %v", c.in, have, c.want)
		}
	}
}

func TestError(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"x", "illegal character"},
	}

	for _, c := range cases {
		_, err := Lex(c.in)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: have: %q, want: %q", c.in, err, c.want)
		}
	}
}
