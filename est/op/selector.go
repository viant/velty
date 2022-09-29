package op

import (
	"github.com/viant/velty/types"
	"github.com/viant/xunsafe"
	"reflect"
)

//Selector represent data selector
type Selector struct {
	ID   string
	Type reflect.Type
	*xunsafe.Field
	Indirect bool
	Parent   *Selector

	Func         *Func
	Slice        *Slice
	Args         []*Operand
	Placeholder  string
	ParentOffset uintptr
}

//NewSelector create a selector
func NewSelector(id, name string, sType reflect.Type, parent *Selector) *Selector {
	xField := newXField(name, sType)
	return &Selector{
		Type:     sType,
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

func SelectorWithField(id string, field *xunsafe.Field, parent *Selector, indirect bool, offset uintptr) *Selector {
	return &Selector{
		Type:         field.Type,
		ID:           id,
		Field:        field,
		Parent:       parent,
		Indirect:     indirect,
		ParentOffset: offset,
	}
}

func FunctionSelector(id string, field *xunsafe.Field, aFunc *Func, parent *Selector) *Selector {
	return &Selector{
		Type:     aFunc.ResultType,
		ID:       id,
		Parent:   parent,
		Indirect: true,
		Func:     aFunc,
		Field:    field,
	}
}

func SliceSelector(id string, placeholder string, sliceOperand, indexOperand *Operand, parent *Selector) (*Selector, error) {
	toInt, err := types.ToInt(indexOperand.Type)
	if err != nil {
		return nil, err
	}

	return &Selector{
		Type:     parent.Type.Elem(),
		ID:       id,
		Indirect: true,
		Parent:   parent,
		Slice: &Slice{
			SliceOperand: sliceOperand,
			IndexOperand: indexOperand,
			ToInter:      toInt,
			XSlice:       xunsafe.NewSlice(parent.Type),
		},
		Placeholder: placeholder,
	}, nil
}
