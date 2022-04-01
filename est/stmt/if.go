package stmt

import (
	est2 "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"unsafe"
)

type If struct {
	ElseIf    est2.Compute
	Block     est2.Compute
	Condition *op2.Operand
}

func (i *If) computeWithoutElse(state *est2.State) unsafe.Pointer {
	if *(*bool)(i.Condition.Exec(state)) {
		return i.Block(state)
	}
	return nil
}

func (i *If) compute(state *est2.State) unsafe.Pointer {
	if *(*bool)(i.Condition.Exec(state)) {
		return i.Block(state)
	}
	return i.ElseIf(state)
}

func NewIf(condition *op2.Expression, block, elseIf est2.New) (est2.New, error) {
	return func(control est2.Control) (est2.Compute, error) {
		result := &If{}
		var err error

		result.Condition, err = condition.Operand(control)
		if err != nil {
			return nil, err
		}

		result.Block, err = block(control)
		if err != nil {
			return nil, err
		}

		if elseIf != nil {
			result.ElseIf, err = elseIf(control)
			if err != nil {
				return nil, err
			}
		}

		if elseIf == nil {
			return result.computeWithoutElse, nil
		}
		return result.compute, nil
	}, nil
}
