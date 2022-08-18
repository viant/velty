package op

import (
	est "github.com/viant/velty/est"
	"reflect"
	"unsafe"
)

type Expression struct {
	LiteralPtr *unsafe.Pointer
	Type       reflect.Type
	*Selector
	est.New
}

func (e *Expression) Operand(control est.Control) (*Operand, error) {
	operand := &Operand{}

	if e.LiteralPtr != nil {
		operand.LiteralPtr = e.LiteralPtr
		operand.SetType(e.Type)
		return operand, nil
	}

	if e.Selector != nil {
		operand.Sel = e.Selector
		operand.SetType(e.Selector.Type)
	}

	if e.Selector != nil && e.Selector.Func != nil {
		operand.SetType(e.Func.ResultType)
		operand.Comp = e.funcCall()
		return operand, nil
	}

	if e.Selector != nil {
		if e.Selector != nil && e.Selector.Indirect {
			operand.Comp = e.newIndirectSelector()
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

func (e *Expression) funcCall() est.Compute {
	upstream := Upstream(e.Selector, false)
	return func(state *est.State) unsafe.Pointer {
		return upstream(state)
	}
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

func (e *Expression) newIndirectSelector() est.Compute {
	upstream := Upstream(e.Selector, true)
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
