package se

var prelude = Scope{symbols: map[string]*Ident{
//"pi":  Const{ValueOf(math.Pi)},
//"add": ReflectFunc(reflect.ValueOf(add)),
//"mul": ReflectFunc(reflect.ValueOf(mul)),
}}

//func add(x, y float64) float64 { return x + y }
//func mul(x, y float64) float64 { return x * y }
//
//type ReflectFunc reflect.Value
//
//func (f ReflectFunc) Apply(args []Value) Value {
//	argv := make([]reflect.Value, len(args))
//	for i, a := range args {
//		argv[i] = reflect.ValueOf(a)
//	}
//	ret := (reflect.Value(f)).Call(argv)
//	if len(ret) != 1 {
//		panic(fmt.Sprintf("cannot handle %v return values", len(ret)))
//	}
//	return ret[0].Interface()
//}
//
//func (f ReflectFunc) Eval() Value               { return Value(f) }
//func (f ReflectFunc) PrintTo(w io.Writer) Value { return Value(f) }
