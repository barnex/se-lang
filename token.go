package se

import (
	"fmt"
)

// A Token represents a textual element like a word, number, ...
type Token struct {
	TType
	Value string
}

func (t Token) String() string {
	switch t.TType {
	default:
		return t.Value
	case TEOF:
		return t.TType.String()
	}
	return fmt.Sprintf("%v", t.Value)
}

// TType is a token type.
type TType int

const (
	tError  TType = iota
	TAdd          // +
	TAssign       // =
	TComma        // ,
	TDiv          // /
	TEOF          // end-of-file
	TEquals       // ==
	TGe           // >=
	TGt           // >
	TIdent        // identifer
	TLBrace       // {
	TLParen       // (
	TLambda       // ->
	TLe           // <=
	TLt           // <
	TMinus        // -
	TMul          // *
	TNum          // number
	TQuote        // '
	TRBrace       // }
	TRParen       // )
	TString       // string
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
	TQuote:  "'",
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
