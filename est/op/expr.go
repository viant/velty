package op

import (
	est "github.com/viant/velty/est"
	"github.com/viant/xunsafe/converter"
	"reflect"
	"unsafe"
)

type ShouldRefLast bool
type ShouldDerefLast bool

type Expression struct {
	LiteralPtr *unsafe.Pointer
	Type       reflect.Type
	*Selector
	est.New
	Unify converter.UnifyFn
}

func (e *Expression) Operand(control est.Control, options ...interface{}) (*Operand, error) {
	var shouldDerefLast bool
	var shouldRefLast bool

	for _, option := range options {
		switch actual := option.(type) {
		case ShouldDerefLast:
			shouldDerefLast = bool(actual)
		case ShouldRefLast:
			shouldRefLast = bool(actual)
		}
	}

	operand := &Operand{}
	operand.SetUnifier(e.Unify)

	if e.Type != nil {
		operand.SetType(e.Type)
	}

	if e.Selector != nil {
		operand.Sel = e.Selector
		operand.trySetType(e.Selector.Type)
	}

	if e.LiteralPtr != nil {
		operand.LiteralPtr = e.LiteralPtr
		return operand, nil
	}

	if e.Selector != nil {
		if e.Selector != nil && e.Selector.Indirect {
			operand.Comp = e.newIndirectSelector(shouldDerefLast, shouldRefLast)
		}

		return operand, nil

	}

	//operand.SetType(e.Type)

	compute, err := e.New(control)
	if err != nil {
		return nil, err
	}
	operand.Comp = compute
	return operand, nil
}

func (e *Expression) funcCall(derefLast bool, refLast bool) est.Compute {
	upstream := Upstream(e.Selector, derefLast, refLast)
	return func(state *est.State) unsafe.Pointer {
		return upstream(state)
	}
}

type Expressions []*Expression

func (e Expressions) Operands(control est.Control, shouldDerefLast bool) ([]*Operand, error) {
	var result = make([]*Operand, len(e))
	var err error
	for i, item := range e {
		if result[i], err = item.Operand(control); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (e *Expression) newIndirectSelector(shouldDerefLast bool, refLast bool) est.Compute {
	upstream := Upstream(e.Selector, shouldDerefLast, refLast)
	return func(state *est.State) unsafe.Pointer {
		ret := upstream(state)
		return ret
	}
}

func NewExpression(selector *Selector) *Expression {
	var litPtr *unsafe.Pointer
	if selector.Literal != nil {
		litPtr = &selector.Literal
	}

	return &Expression{
		Selector:   selector,
		LiteralPtr: litPtr,
	}
}
