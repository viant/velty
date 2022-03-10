package stmt

import (
	"github.com/viant/velty/est"
	"unsafe"
)

type block struct {
	stmt []est.Compute
}

func (s *block) compute(state *est.State) unsafe.Pointer {
	var result unsafe.Pointer
	for i := 0; i < len(s.stmt); i++ {
		result = s.stmt[i](state)
	}
	return result
}

type stmt1 struct {
	compute est.Compute
}

type stmt2 struct {
	stmt1
	est.Compute
}

func (s *stmt2) compute(state *est.State) unsafe.Pointer {
	s.stmt1.compute(state)
	return s.Compute(state)
}

func new2Stmt(args []est.Compute) *stmt2 {
	return &stmt2{
		stmt1:   stmt1{compute: args[0]},
		Compute: args[1],
	}
}

func nop(mem *est.State) unsafe.Pointer {
	return nil
}

func NewBlock(stmtsNew []est.New) est.New {
	computers := est.Computers(stmtsNew)
	return func(control est.Control) (est.Compute, error) {
		stmts, err := computers.New(control)
		if err != nil {
			return nil, err
		}
		switch len(stmts) {
		case 0:
			return nop, nil
		case 1:
			return stmts[0], nil
		case 2:
			return new2Stmt(stmts).compute, nil
		default:
			b := &block{stmt: stmts}
			return b.compute, nil
		}
	}
}
