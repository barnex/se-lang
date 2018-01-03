package e

import (
	"reflect"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		src  string
		want Node
	}{
		{`1+1`, num(2)},
		{`1+2*3`, num(7)},
		{`1*2+3`, num(5)},
	}

	for _, c := range cases {
		n, err := Parse(strings.NewReader(c.src))
		if err != nil {
			t.Errorf("%v: %v", c.src, err)
			continue
		}
		have, err := EvalNode(n)
		if err != nil {
			t.Errorf("%v: %v", c.src, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%v: have %v, want %v", c.src, have, c.want)
		}
	}
}
