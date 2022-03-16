package stmt

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"unsafe"
)

type For struct {
	Block     est.Compute
	Init      est.Compute
	Condition *op.Operand

	Post est.Compute
}

func (f *For) compute(state *est.State) unsafe.Pointer {
	var ptr unsafe.Pointer
	f.Init(state)
	for *(*bool)(f.Condition.Exec(state)) {
		ptr = f.Block(state)
		f.Post(state)
	}

	return ptr
}

func ForLoop(init, post est.New, condition *op.Expression, block est.Compute) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		forLoop := &For{}
		var err error

		forLoop.Condition, err = condition.Operand(control)
		if err != nil {
			return nil, err
		}

		forLoop.Init, err = init(control)
		if err != nil {
			return nil, err
		}

		forLoop.Post, err = post(control)
		if err != nil {
			return nil, err
		}

		forLoop.Block = block

		return forLoop.compute, nil
	}, nil
}
