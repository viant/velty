package stmt

import (
	"github.com/viant/velty/internal/est"
	"github.com/viant/velty/internal/est/op"
	"unsafe"
)

type If struct {
	ElseIf    est.Compute
	Block     est.Compute
	Condition *op.Operand
}

func (i *If) computeWithoutElse(state *est.State) unsafe.Pointer {
	if *(*bool)(i.Condition.Exec(state)) {
		return i.Block(state)
	}
	return nil
}

func (i *If) compute(state *est.State) unsafe.Pointer {
	if *(*bool)(i.Condition.Exec(state)) {
		return i.Block(state)
	}
	return i.ElseIf(state)
}

func NewIf(condition *op.Expression, block, elseIf est.New) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
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
