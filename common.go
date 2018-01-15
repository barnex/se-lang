/*
package se contains a few utilities shared among the other packages.

The real functionality is implemented in these packages:

	lex: Lexical scanner for SE source text.
	ast: Parser and Abstract Syntax Tree
	typ: Typechekcer
	std: Standard library
	eva: Intermediate Representation & evaluator
	jit: Just-In-Time compiler
*/
package se

import (
	"fmt"
	"text/scanner"
)

type Position struct {
	scanner.Position
}

type Error struct {
	Msg string
}

func Errorf(format string, x ...interface{}) Error {
	return Error{Msg: fmt.Sprintf(format, x...)}
}

func (e Error) Error() string {
	return e.Msg
}

func IsSEError(e error) bool {
	_, ok := e.(Error)
	return ok
}
