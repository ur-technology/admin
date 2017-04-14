package errors

import "fmt"

type base struct {
	s string
}

func (b base) Error() string {
	return b.s
}

type Internal struct{ base }

func NewInternal(desc string, args ...interface{}) (e Internal) {
	e.s = fmt.Sprintf(desc, args...)
	return e
}

type External struct{ base }

func NewExternal(desc string, args ...interface{}) (e External) {
	e.s = fmt.Sprintf(desc, args...)
	return e
}

type User struct{ base }

func NewUser(desc string, args ...interface{}) (e User) {
	e.s = fmt.Sprintf(desc, args...)
	return e
}
