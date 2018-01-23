package eva

import (
	"fmt"
	"strings"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		src  string
		want interface{}
	}{
		// arithmetic
		{`1`, 1},
		{`1+2`, 3},
		{`1+2+3+4`, 10},
		{`(1+2)+(3+4)`, 10},
		{`1*2*3*4`, 24},
		{`(1*2)*(3*4)`, 24},
		{`-1`, -1},
		{`2-1`, 1},

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
		{`1+2*3%4`, 3},

		// cond
		{`true? 1 : 2`, 1},
		{`false? 1 : 2`, 2},
		{`((x,y)->1+2==3? x : y)(111,222)`, 111},
		{`((x,y)->x>y?x:y)(111,222)`, 222}, // max
		{`((x,y)->x<y?x:y)(111,222)`, 111}, // min
		{`{x=111; y=222; x>y?x:y}`, 222},   // max, inlined

		// lambda
		{`(x->x)(1)`, 1},                         // identity function
		{`(()->7)()`, 7},                         // constant function
		{`(x->-x)(1)`, -1},                       // lambda: negative
		{`(x->x*x)(3)`, 9},                       // lambda: square
		{`(x->x+x)(3)`, 6},                       // lambda: double
		{`((x,y)->x+y)(1,2)`, 3},                 // lambda: sum
		{`((x,y)->x)(1,2)`, 1},                   // lambda: first
		{`((x,y)->y)(1,2)`, 2},                   // lambda: second
		{`((f,i)->f(i))((x->x*x), 3)`, 9},        // lambda: apply f to i
		{`( (f,i)->f(f(i)) ) ( (x->x*2), 1)`, 4}, // lambda: apply f twice

		// closure
		{`(x->()->x)(1)()`, 1},
		{`(x->y->x+y)(1)(2)`, 3},         // close over parent
		{`(x->y->z->x+y+z)(1)(2)(3)`, 6}, // transitive

		// block, assign
		{`(()->{x=1;x})()`, 1},
		{`(()->{x=1;y=x+2;x+y})()`, 4},
		{`f=(x,y)->x+y; f(1,2)`, 3},
		{`(()->{f=(x,y)->x+y; f(1,2)})()`, 3},
		{`{x=1;x}`, 1},
		{`{{{x=1;x}}}`, 1},
		{`{x=1;y=x+2;x+y}`, 4},
		{`{f=(x,y)->x+y; f(1,2)}`, 3},

		// program
		{`x=1;x`, 1},
		{`max=(x,y)->x>y?x:y; max(1,2)`, 2},
		{`id=x->x; id(2)`, 2},
		{`max=(x,y)->x>y?x:y; max(2,1)`, 2},

		// weird
		//{`{add}(1,2)`, 3},
		//{`{f=add;f}(1,2)`, 3},

		// recursion
		{`fac=(n)->{n <= 1? n: n*fac(n-1)}; fac(6)`, 720},
		{`fac=(n)->(n <= 1? n: n*fac(n-1)); fac(6)`, 720},
		{`fib=(n)->(n<=2)?1:(fib(n-1)+fib(n-2)); fib(12)`, 144},
	}

	for _, c := range cases {
		fmt.Println("\n", c.src)
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
