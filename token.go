package main

// A Token represents a textual element like a word, number, ...
type Token struct {
	Type
	Value string
}

type Type int

// All possible token types
const (
	tErr Type = iota // error, internal use only
	TNum             // number
	TEOF             // end-of-file
)

func (t Token) String() string {
	return tokenName[t.Type] + "(" + t.Value + ")"
}

var tokenName = map[Type]string{
	tErr: "Err",
	TNum: "Num",
	TEOF: "EOF",
}
