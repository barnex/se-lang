package main

import (
	"reflect"
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
		{" 45 678 ", []Token{{TNum, "45"}, {TNum, "678"}, {TEOF, ""}}},
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
