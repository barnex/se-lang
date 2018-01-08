package se

import (
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	cases := []struct {
		src  string
		want interface{}
	}{
		{`1`, 1.0},
		{`1+1`, 2.0},
		{`1+2+3+4`, 10.0},
		{`(1+2)+(3+4)`, 10.0},
		{`1*2*3*4`, 24.0},
		{`(1*2)*(3*4)`, 24.0},
		{`(x->x)(1)`, 1.0},                         // identity function
		{`(x->x*x)(3)`, 9.0},                       // lambda: square
		{`((x,y)->x+y)(1,2)`, 3.0},                 // lambda: sum
		{`((f,i)->f(i))((x->x*x), 3)`, 9.0},        // lambda: apply f to i
		{`( (f,i)->f(f(i)) ) ( (x->x*2), 1)`, 4.0}, // lambda: apply f twice
		//{`((f,i)->f(i))(x->x*x, 3)`, 9.0},
		//{`d->x->x+d`, 1.0},
	}

	for _, c := range cases {
		prog, err := Compile(strings.NewReader(c.src))
		if err != nil || prog == nil {
			t.Fatalf("%v: have %v, %v", c.src, prog, err)
		}
		if have := prog.Eval(); have != c.want {
			t.Errorf("%v: have %v, want: %v", c.src, have, c.want)
		}
	}
}
