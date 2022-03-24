package expr

import (
	"reflect"
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

//StringExpression creates string literal
func StringExpression(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(""),
	}
}

//NumberExpression creates number literal
func NumberExpression(value string) *Literal {
	numType := reflect.TypeOf(0.0)
	if !(strings.Contains(value, ".") || strings.Contains(value, "e")) {
		numType = reflect.TypeOf(0)
	}
	return &Literal{
		Value: value,
		RType: numType,
	}
}

//BoolExpression creates bool literal
func BoolExpression(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(true),
	}
}
