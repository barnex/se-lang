package main

type SyntaxError struct {
	msg string
	// TODO: pos
}

func (e SyntaxError) Error() string {
	return e.msg
}
