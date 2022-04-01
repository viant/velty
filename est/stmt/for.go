package stmt

import (
	est2 "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"unsafe"
)

type For struct {
	Block     est2.Compute
	Init      est2.Compute
	Condition *op2.Operand

	Post est2.Compute
}

func (f *For) compute(state *est2.State) unsafe.Pointer {
	var ptr unsafe.Pointer
	f.Init(state)
	for *(*bool)(f.Condition.Exec(state)) {
		ptr = f.Block(state)
		f.Post(state)
	}

	return ptr
}

func ForLoop(init, post est2.New, condition *op2.Expression, block est2.Compute) (est2.New, error) {
	return func(control est2.Control) (est2.Compute, error) {
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
