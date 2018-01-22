se-lang is a toy functional programming language. Intended for fun with lambda expressions, not any practical purpose.

Examples:

## Arithmetic
```
(1+2+3)*(3+4)     // 42
```

## Comparison
```
2 < 3            // true
2 >= 3           // false
2 == 3           // false
2 != 3           // true
```

## Conditional
```
x>0? x: 0        // if x>0 then x else 0
max = x>y? x: y
```

## Lambda's

```
x -> x*x     // a function that computes the square
(x->x*x)(3)  // 9
```

```
max = (x,y) -> x>y? x: y
max(1,2)                  // 2
```

```
(f,a) -> f(a)                  // a function that applies a function to an argument
(f,a) -> f(f(a))               // a function that applies a function twice
((f,a)->f(f(a))) ((x->x*x), 3) // apply square to 3, twice (result: 81)
```
