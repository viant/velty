package op

import (
	"github.com/viant/velty/est"
	"reflect"
	"unsafe"
)

type Expression struct {
	LiteralPtr *unsafe.Pointer
	Type       reflect.Type
	*est.Selector
	est.New
}

func (e *Expression) Operand(control est.Control) (*Operand, error) {
	operand := &Operand{}

	if e.LiteralPtr != nil {
		operand.LiteralPtr = e.LiteralPtr
		operand.Type = e.Type
		return operand, nil
	}
	if e.Selector != nil {
		//TODO check direct (no ptr, slice etc ...)
		operand.Offset = &e.Selector.Offset
		operand.Type = e.Selector.Type
		operand.Sel = e.Selector

		return operand, nil

		//		operand.Sel = e.selectorExists
		//		return operand, nil
	}
	compute, err := e.New(control)
	if err != nil {
		return nil, err
	}
	operand.Comp = compute
	return operand, nil
}

type Expressions []*Expression

func (e Expressions) Operands(control est.Control) ([]*Operand, error) {
	var result = make([]*Operand, len(e))
	var err error
	for i, item := range e {
		if result[i], err = item.Operand(control); err != nil {
			return nil, err
		}
	}
	return result, nil
}
