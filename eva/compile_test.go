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
		{`-1`, -1.0},
		{`2-1`, 1.0},

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
		{`!true`, false},
		{`!false`, true},

		// precedence
		{`true==false||false==false`, true},
		{`1+1==2&&3<4`, true},
		//{`1+2*3%4`, 3},

		// cond
		{`true? 1 : 2`, 1.0},
		{`false? 1 : 2`, 2.0},
		{`((x,y)->1+2==3? x : y)(111,222)`, 111.0},
		{`((x,y)->x>y?x:y)(111,222)`, 222.0}, // max
		{`((x,y)->x<y?x:y)(111,222)`, 111.0}, // min
		//{`{x=111; y=222; x>y?x:y}`, 222.0},   // max, inlined

		// lambda
		{`(x->x)(1)`, 1.0},                         // identity function
		{`(()->7)()`, 7.0},                         // constant function
		{`(x->-x)(1)`, -1.0},                       // lambda: negative
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

		// block, assign
		{`(()->{x=1;x})()`, 1.0},
		{`(()->{x=1;y=x+2;x+y})()`, 4.0},
		{`(()->{f=(x,y)->x+y; f(1,2)})()`, 3.0},
		{`{x=1;x}`, 1.0},
		{`{{{x=1;x}}}`, 1.0},
		{`{x=1;y=x+2;x+y}`, 4.0},
		{`{f=(x,y)->x+y; f(1,2)}`, 3.0},
		//{`(()->{f=()->f(); f()})()`, 3.0},

		// program
		{`x=1;x`, 1.0},
		{`max=(x,y)->x>y?x:y; max(1,2)`, 2.0},
		{`max=(x,y)->x>y?x:y; max(2,1)`, 2.0},
	}

	for _, c := range cases {
		prog, err := Compile(strings.NewReader(c.src))
		if err != nil || prog == nil {
			t.Errorf("%v: have %v, %v", c.src, prog, err)
			continue
		}
		if have, err := Eval(prog); err != nil || have != c.want {
			t.Errorf("%v: have %v, %v, want: %v", c.src, have, err, c.want)
		}
	}
}
