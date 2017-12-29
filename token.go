package e

import (
	"fmt"
)

// A Token represents a textual element like a word, number, ...
type Token struct {
	TType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%v", t.Value)
}

// TType is a token type.
// Positive values are just unicode runes,
// Negative values are defined below.
type TType int

const (
	tError  = iota
	TAdd    // +
	TAssign // =
	TComma  // ,
	TDiv    // /
	TEOF    // end-of-file
	TEquals // ==
	TGe     // >=
	TGt     // >
	TIdent  // identifer
	TLBrace // {
	TLParen // (
	TLambda // ->
	TLe     // <=
	TLt     // <
	TMinus  // -
	TMul    // *
	TNum    // number
	TRBrace // }
	TRParen // )
	TString // string
)

var ttypeString = map[TType]string{
	TAdd:    "+",
	TAssign: "=",
	TComma:  ",",
	TDiv:    "/",
	TEOF:    "EOF",
	TEquals: "==",
	TGe:     ">=",
	TGt:     ">",
	TIdent:  "identifer",
	TLBrace: "{",
	TLParen: "(",
	TLambda: "->",
	TLe:     "<=",
	TLt:     "<",
	TMinus:  "-",
	TMul:    "*",
	TNum:    "number",
	TRBrace: "}",
	TRParen: ")",
	TString: "string",
}

func (t TType) String() string {
	if str, ok := ttypeString[t]; ok {
		return str
	}
	panic(fmt.Sprintf("bad TType: %q", rune(t)))
}
