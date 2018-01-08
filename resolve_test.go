package se

//func TestResolveFound(t *testing.T) {
//	cases := []string{
//		`1`,
//		`(x)->(x)`,
//		`(x)->(y->x)`,
//	}
//
//	for _, src := range cases {
//		n, err := Parse(strings.NewReader(src))
//		if err != nil {
//			t.Fatal(err)
//		}
//		if err := resolve(n); err != nil {
//			t.Errorf("%v: %v", src, err)
//		}
//	}
//}
//
//func TestResolveNotFound(t *testing.T) {
//	cases := []string{
//		`x`,
//		`(x)->(y)`,
//		`(x)->(x->y)`,
//	}
//
//	for _, src := range cases {
//		n, err := Parse(strings.NewReader(src))
//		if err != nil {
//			t.Fatal(err)
//		}
//		if err := resolve(n); err == nil {
//			t.Errorf("%v: expected error", src)
//		}
//	}
//
//}
//
//func resolve(n Node) (e error) {
//	defer func() {
//		switch p := recover().(type) {
//		case nil:
//		default:
//			panic(p)
//		case *SyntaxError:
//			e = p
//		}
//	}()
//	//Resolve(&prelude, n)
//	return nil
//}
