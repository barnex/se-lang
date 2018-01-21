package lex

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
	TError   TType = iota
	TAdd           // +
	TAnd           // &&
	TAssign        // =
	TComma         // ,
	TDiv           // /
	TEOF           // end-of-file
	TEq            // ==
	TGe            // >=
	TGt            // >
	TIdent         // identifer
	TLBrace        // {
	TLParen        // (
	TLambda        // ->
	TLe            // <=
	TLt            // <
	TMinus         // -
	TMod           // %
	TMul           // *
	TNEq           // !=
	TNot           // !
	TNum           // number
	TOr            // ||
	TQuote         // '
	TRBrace        // }
	TRParen        // )
	TSemicol       // ;
	TString        // string
)

var ttypeString = map[TType]string{
	TAdd:     "+",
	TAnd:     "&&",
	TAssign:  "=",
	TComma:   ",",
	TDiv:     "/",
	TEOF:     "EOF",
	TEq:      "==",
	TGe:      ">=",
	TGt:      ">",
	TIdent:   "identifer",
	TLBrace:  "{",
	TLParen:  "(",
	TLambda:  "->",
	TLe:      "<=",
	TLt:      "<",
	TMinus:   "-",
	TMod:     "%",
	TMul:     "*",
	TNEq:     "!=",
	TNot:     "!",
	TNum:     "number",
	TOr:      "||",
	TQuote:   "'",
	TRBrace:  "}",
	TRParen:  ")",
	TSemicol: ";",
	TString:  "string",
}

func (t TType) String() string {
	if str, ok := ttypeString[t]; ok {
		return str
	}
	panic(fmt.Sprintf("bad TType: %q", rune(t)))
}
