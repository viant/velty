package est

import (
	"github.com/viant/velty/ast"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type Args []Compute

func (a Args) ToPtrs(ret unsafe.Pointer, state *State) []unsafe.Pointer {
	result := make([]unsafe.Pointer, len(a)+1)
	result[0] = ret
	for i, _ := range a {
		result[i+1] = a[i](state)
	}
	return result
}

//Selector represent data selector
type Selector struct {
	ID string
	*xunsafe.Field
	Indirect bool
	Parent   *Selector

	rType         reflect.Type
	Func          *Func
	FuncArguments []ast.Expression
	Args          []Compute
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

func (s *Selector) Type() reflect.Type {
	if s.rType != nil {
		return s.rType
	}

	return s.Field.Type
}

func (s *Selector) SetType(rType reflect.Type) {
	s.rType = rType
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

func SelectorWithField(id string, field *xunsafe.Field, parent *Selector) *Selector {
	isIndirectParent := false
	if parent != nil {
		isIndirectParent = parent.Indirect
	}
	return &Selector{
		ID:       id,
		Field:    field,
		Parent:   parent,
		Indirect: isIndirectParent || field.Kind() == reflect.Ptr || field.Kind() == reflect.Slice,
	}
}

func FunctionSelector(id string, aFunc *Func, args []ast.Expression, parent *Selector) *Selector {
	return &Selector{
		ID:            id,
		Parent:        parent,
		Indirect:      true,
		FuncArguments: args,
		Func:          aFunc,
		rType:         aFunc.ResultType,
	}
}
