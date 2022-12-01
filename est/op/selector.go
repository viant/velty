package op

import (
	"github.com/viant/xunsafe"
	types "github.com/viant/xunsafe/converter"
	"reflect"
)

//Selector represent data selector
type Selector struct {
	ID   string
	Type reflect.Type
	*xunsafe.Field
	Indirect bool
	Parent   *Selector

	Func          *Func
	Slice         *Slice
	Args          []*Operand
	Placeholder   string
	ParentOffset  uintptr
	Map           *Map
	InterfaceExec *Interface
	Cycle         *Selector
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
	return newSelector(id, field, parent, indirect, offset, nil)
}

func NewCycleSelector(id string, field *xunsafe.Field, parent *Selector, indirect bool, offset uintptr, cycleSelector *Selector) *Selector {
	return newSelector(id, field, parent, indirect, offset, cycleSelector)
}

func newSelector(id string, field *xunsafe.Field, parent *Selector, indirect bool, offset uintptr, cycleSelector *Selector) *Selector {
	return &Selector{
		Type:         field.Type,
		ID:           id,
		Field:        field,
		Parent:       parent,
		Indirect:     indirect,
		ParentOffset: offset,
		Cycle:        cycleSelector,
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
	slice, err := newSlice(sliceOperand, indexOperand, parent)
	if err != nil {
		return nil, err
	}

	return &Selector{
		Type:        parent.Type.Elem(),
		ID:          id,
		Indirect:    true,
		Parent:      parent,
		Slice:       slice,
		Placeholder: placeholder,
	}, nil
}

func newSlice(sliceOperand *Operand, indexOperand *Operand, parent *Selector) (*Slice, error) {
	toInt, err := types.Unify(indexOperand.Type, intType)
	if err != nil {
		return nil, err
	}

	return &Slice{
		SliceOperand: sliceOperand,
		IndexOperand: indexOperand,
		ToInter:      toInt.Y,
		XSlice:       xunsafe.NewSlice(parent.Type),
	}, nil
}

func NewMapSelector(id string, placeholder string, mapOperand, indexOperand *Operand, parent *Selector) (*Selector, error) {
	return &Selector{
		Type:        parent.Type.Elem(),
		ID:          id,
		Indirect:    true,
		Parent:      parent,
		Placeholder: placeholder,
		Map:         newMap(mapOperand, indexOperand, parent),
	}, nil
}

func NewInterfaceSelector(id string, placeholder string, xOperand, indexOperand *Operand, parent *Selector) (*Selector, error) {
	aSlice, err := newSlice(xOperand, indexOperand, parent)
	if err != nil {
		return nil, err
	}

	return &Selector{
		Type:        parent.Type,
		ID:          id,
		Indirect:    true,
		Parent:      parent,
		Placeholder: placeholder,
		InterfaceExec: &Interface{
			xOperand: xOperand,
			aMap:     newMap(xOperand, indexOperand, parent),
			aSlice:   aSlice,
		},
	}, nil
}

func newMap(mapOperand *Operand, indexOperand *Operand, parent *Selector) *Map {
	rType := mapOperand.Type
	switch rType.Kind() {
	case reflect.Map:
		rType = rType.Elem()
	}

	elemKind := rType.Kind()
	return &Map{
		mapOperand:   mapOperand,
		indexOperand: indexOperand,
		isValueIface: elemKind == reflect.Interface,
		elemKind:     elemKind,
	}
}
