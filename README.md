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
((f,a)->f(f(a))) ((x->x*x), 3) // apply square to 3, twice (result: 81)
```
