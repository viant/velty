package expr

import (
	"reflect"
	"strconv"
	"strings"
)

//Literal represents constant value i.e. 1, "foo" etc
type Literal struct {
	Value string
	RType reflect.Type
}

func (l *Literal) Type() reflect.Type {
	return l.RType
}

//StringLiteral creates string literal
func StringLiteral(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(""),
	}
}

//NumberLiteral creates number literal
func NumberLiteral(value string) *Literal {
	numType := reflect.TypeOf(0.0)
	if !(strings.Contains(value, ".") || strings.Contains(value, "e")) {
		numType = reflect.TypeOf(0)
	}
	return &Literal{
		Value: value,
		RType: numType,
	}
}

func Number(value int) *Literal {
	numType := reflect.TypeOf(value)
	return &Literal{
		Value: strconv.Itoa(value),
		RType: numType,
	}
}

//BoolLiteral creates bool literal
func BoolLiteral(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(true),
	}
}
