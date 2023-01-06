package op

import (
	est "github.com/viant/velty/est"
	"github.com/viant/xunsafe/converter"
	"reflect"
	"unsafe"
)

type Expression struct {
	LiteralPtr *unsafe.Pointer
	Type       reflect.Type
	*Selector
	est.New
	Unify converter.UnifyFn
}

func (e *Expression) Operand(control est.Control, shouldDerefLast bool) (*Operand, error) {
	var unify func(pointer unsafe.Pointer) unsafe.Pointer
	if e.Unify != nil {
		unify = func(pointer unsafe.Pointer) unsafe.Pointer {
			ptr, _ := e.Unify(pointer)
			return ptr
		}
	}
	operand := &Operand{
		unify: unify,
	}

	if e.Type != nil {
		operand.SetType(e.Type)
	}

	if e.LiteralPtr != nil {
		operand.LiteralPtr = e.LiteralPtr
		return operand, nil
	}

	if e.Selector != nil {
		operand.Sel = e.Selector
		operand.trySetType(e.Selector.Type)
	}

	if e.Selector != nil && e.Selector.Func != nil {
		operand.trySetType(e.Func.ResultType)
		operand.Comp = e.funcCall(shouldDerefLast)
		return operand, nil
	}

	//if e.Selector != nil && e.Selector.Slice != nil {
	//	operand.SetType(e.Type)
	//}

	if e.Selector != nil {
		if e.Selector != nil && e.Selector.Indirect {
			operand.Comp = e.newIndirectSelector(shouldDerefLast)
		}

		return operand, nil

	}

	operand.SetType(e.Type)
	compute, err := e.New(control)
	if err != nil {
		return nil, err
	}
	operand.Comp = compute
	return operand, nil
}

func (e *Expression) funcCall(derefLast bool) est.Compute {
	upstream := Upstream(e.Selector, derefLast)
	return func(state *est.State) unsafe.Pointer {
		return upstream(state)
	}
}

type Expressions []*Expression

func (e Expressions) Operands(control est.Control, shouldDerefLast bool) ([]*Operand, error) {
	var result = make([]*Operand, len(e))
	var err error
	for i, item := range e {
		if result[i], err = item.Operand(control, shouldDerefLast); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (e *Expression) newIndirectSelector(shouldDerefLast bool) est.Compute {
	upstream := Upstream(e.Selector, shouldDerefLast)
	return func(state *est.State) unsafe.Pointer {
		ret := upstream(state)
		return ret
	}
}

func NewExpression(selector *Selector) *Expression {
	return &Expression{
		Selector: selector,
	}
}
