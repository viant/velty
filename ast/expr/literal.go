package expr

import "reflect"

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
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(0.0),
	}
}

func BoolExpression(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(true),
	}
}

func FloatExpression(value string) *Literal {
	return &Literal{
		Value: value,
		RType: reflect.TypeOf(0.0),
	}
}
