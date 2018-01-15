package se

import "fmt"

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
