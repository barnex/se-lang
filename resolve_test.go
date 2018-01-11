package se

import (
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	cases := []struct {
		src string
	}{
		{`1`},
		{`1+1`},
		{`1+2+3+4`},
		{`(1+2)+(3+4)`},
		{`1*2*3*4`},
		{`(1*2)*(3*4)`},
		{`(x->x)(1)`},
		{`(x->x*x)(3)`},
		{`((x,y)->x+y)(1,2)`},
		{`((f,i)->f(i))((x->x*x), 3)`},
		{`( (f,i)->f(f(i)) ) ( (x->x*2), 1)`},
		{`(x->y->x+y)(1)(2)`},
	}

	for _, c := range cases {
		ast, err := Parse(strings.NewReader(c.src))
		if err != nil || ast == nil {
			t.Fatalf("%v: have %v, %v", c.src, ast, err)
		}
		Resolve(ast)
	}
}
