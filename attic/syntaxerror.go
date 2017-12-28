package main

type SyntaxError struct {
	msg string
	// TODO: pos
}

func (e SyntaxError) Error() string {
	return e.msg
}

func IsSyntaxError(e error) bool {
	_, ok := e.(SyntaxError)
	return ok
}
