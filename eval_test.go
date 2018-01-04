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
	//{`1`, num(1)},
	//{`pi`, num(math.Pi)},
	//{`1+1`, num(2)},
	//{`1+2*3`, num(7)},
	//{`1*2+3`, num(5)},
	}

	for _, c := range cases {
		n, err := Parse(strings.NewReader(c.src))
		if err != nil {
			t.Errorf("%v: %v", c.src, err)
			continue
		}
		have, err := EvalSafe(n)
		if err != nil {
			t.Errorf("%v: %v", c.src, err)
			continue
		}
		if !reflect.DeepEqual(have, c.want) {
			t.Errorf("%v: have %v, want %v", c.src, have, c.want)
		}
	}
}
