package e

import "fmt"

type SyntaxError struct {
	Msg string
	Position
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("line %v:%v: %v", e.Line, e.Column, e.Msg)
}

func IsSyntaxError(e error) bool {
	_, ok := e.(*SyntaxError)
	return ok
}
