package main

import (
	"reflect"
	"testing"
)

func TestParseExpr(t *testing.T) {

	cases := []struct {
		in   string
		want Expr
	}{
		{`1`, Num{1}},
		{` 1 `, Num{1}},
		{`(1)`, Num{1}},
		{` ( 1 ) `, Num{1}},
	}

	for _, c := range cases {
		have, err := Parse(c.in)
		if err != nil {
			t.Errorf("%v: error: %v", c.in, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%v: have %v, want %v", c.in, have, c.want)
		}
	}
}

func TestParseError(t *testing.T) {
	cases := []string{
		`(1`,
		`1)`,
		` ( 1 `,
		` 1 ) `,
	}

	for _, c := range cases {
		_, err := Parse(c)
		if err == nil {
			t.Errorf("%v: expected error", c)
			continue
		}
	}
}
