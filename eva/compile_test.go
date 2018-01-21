package eva

import (
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		src  string
		want interface{}
	}{
		// arithmetic
		{`1`, 1.0},
		{`1+2`, 3.0},
		{`1+2+3+4`, 10.0},
		{`(1+2)+(3+4)`, 10.0},
		{`1*2*3*4`, 24.0},
		{`(1*2)*(3*4)`, 24.0},

		// comparison
		{`1==1`, true},
		{`1==2`, false},
		{`1!=1`, false},
		{`1!=2`, true},
		{`1<2`, true},
		{`2<1`, false},
		{`1>2`, false},
		{`2>1`, true},
		{`1>=2`, false},
		{`1>=1`, true},
		{`2>=1`, true},
		{`1<=2`, true},
		{`1<=1`, true},
		{`2<=1`, false},

		// boolean
		{`true`, true},
		{`false`, false},
		{`true && false`, false},
		{`true && true`, true},
		{`true || false`, true},
		{`true || true`, true},
		//{`!true`, false},
		//{`!false`, true},

		// precedence
		{`true==false||false==false`, true},
		{`1+1==2&&3<4`, true},
		//{`1+2*3%4`, 3},

		// lambda
		{`(x->x)(1)`, 1.0},                         // identity function
		{`(()->7)()`, 7.0},                         // constant function
		{`(x->x*x)(3)`, 9.0},                       // lambda: square
		{`(x->x+x)(3)`, 6.0},                       // lambda: double
		{`((x,y)->x+y)(1,2)`, 3.0},                 // lambda: sum
		{`((x,y)->x)(1,2)`, 1.0},                   // lambda: first
		{`((x,y)->y)(1,2)`, 2.0},                   // lambda: second
		{`((f,i)->f(i))((x->x*x), 3)`, 9.0},        // lambda: apply f to i
		{`( (f,i)->f(f(i)) ) ( (x->x*2), 1)`, 4.0}, // lambda: apply f twice

		// closure
		{`(x->()->x)(1)()`, 1.0},
		{`(x->y->x+y)(1)(2)`, 3.0},         // close over parent
		{`(x->y->z->x+y+z)(1)(2)(3)`, 6.0}, // transitive
	}

	for _, c := range cases {
		prog, err := Compile(strings.NewReader(c.src))
		if err != nil || prog == nil {
			t.Fatalf("%v: have %v, %v", c.src, prog, err)
		}
		if have, err := Eval(prog); err != nil || have != c.want {
			t.Errorf("%v: have %v, %v, want: %v", c.src, have, err, c.want)
		}
	}
}
