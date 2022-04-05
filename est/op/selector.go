package op

import (
	"github.com/viant/velty/internal/ast"
	"github.com/viant/xunsafe"
	"reflect"
)

//Selector represent data selector
type Selector struct {
	ID string
	*xunsafe.Field
	Indirect bool
	Parent   *Selector

	Func          *Func
	FuncArguments []ast.Expression
	Args          []*Operand
}

//NewSelector create a selector
func NewSelector(id, name string, sType reflect.Type, parent *Selector) *Selector {
	xField := newXField(name, sType)
	return &Selector{
		ID:       id,
		Field:    xField,
		Parent:   parent,
		Indirect: sType != nil && (sType.Kind() == reflect.Ptr || sType.Kind() == reflect.Slice),
	}
}

func newXField(name string, sType reflect.Type) *xunsafe.Field {
	field := reflect.StructField{Name: name, Type: sType}
	var xField *xunsafe.Field
	if field.Type != nil && field.Type.Kind() != reflect.Invalid {
		xField = xunsafe.NewField(field)
	} else {
		xField = &xunsafe.Field{Name: name, Type: sType}
	}
	return xField
}

func SelectorWithField(id string, field *xunsafe.Field, parent *Selector, indirect bool) *Selector {
	return &Selector{
		ID:       id,
		Field:    field,
		Parent:   parent,
		Indirect: indirect,
	}
}

func FunctionSelector(id string, field reflect.StructField, aFunc *Func, args []ast.Expression, parent *Selector) *Selector {
	return &Selector{
		ID:            id,
		Parent:        parent,
		Indirect:      true,
		FuncArguments: args,
		Func:          aFunc,
		Field:         xunsafe.NewField(field),
	}
}
