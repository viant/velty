package plan

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/parser"
	"github.com/viant/xunsafe"
	"reflect"
)

func (p *Planner) Compile(template []byte) (*est.Execution, func() *est.State, error) {
	root, err := parser.Parse(template)
	if err != nil {
		return nil, nil, err
	}
	computeNew, err := p.compileBlock(root)
	if err != nil {
		return nil, nil, err
	}
	compute, err := computeNew(*p.Control)
	if err != nil {
		return nil, nil, err
	}
	exec := est.NewExecution(compute)
	newState := func() *est.State {
		mem := reflect.New(p.Type.Type).Interface()
		state := &est.State{
			Mem:       mem,
			MemPtr:    xunsafe.AsPointer(mem),
			Buffer:    est.NewBuffer(p.bufferSize),
			Selectors: *p.selectors,
			Index:     p.index,
		}
		return state
	}
	return exec, newState, nil
}
