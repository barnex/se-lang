se-lang is a toy functional programming language. Intended for fun with lambda expressions, not any practical purpose.

### Language examples:
```
fac = n ->
	n==1? 1: n*fac(n-1);

fac(6)  // 720
```

```
(x -> x*x)(3)  // square of 3
```

```
max = (x,y) -> x>y? x: y;
max(1,2)                  // 2
```

```
square = x -> x*x;
twice = f -> (x -> f(f(x))); // applies a function twice
(twice(square)) (3)          // 81
```

```
((f,a)->f(f(a))) ((x->x*x), 3) // same as above
```
