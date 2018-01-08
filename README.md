se-lang is a toy functional programming language.

Examples:

```
x -> x*x     // a function that computes the square
(x->x*x)(3)  // apply square to 3 (result: 9)

(f,a) -> f(a)                  // a function that applies a function to an argument
(f,a) -> f(f(a))               // a function that applies a function twice
((f,a)->f(f(a))) ((x->x*x), 3) // apply square to 3, twice (result: 81)
```
