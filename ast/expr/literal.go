package expr

import (
	"reflect"
	"strings"
)

type Literal struct {
	Value string
	RType reflect.Type
}

func (l *Literal) Type() reflect.Type {
	return l.RType
}

func StringExpression(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(""),
	}
}

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

func BoolExpression(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(true),
	}
}
