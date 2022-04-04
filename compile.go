package velty

import (
	est2 "github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast/stmt"
	"github.com/viant/velty/internal/parser"
	"github.com/viant/xunsafe"
	"reflect"
)

//Compile create Execution Plan and State provider for the Execution Plan.
func (p *Planner) Compile(template []byte) (*est2.Execution, func() *est2.State, error) {
	root, err := parser.Parse(template)
	if err != nil {
		return nil, nil, err
	}

	exec, err := p.newExecution(root)
	if err != nil {
		return nil, nil, err
	}

	newState := p.stateProvider()
	return exec, newState, nil
}

func (p *Planner) stateProvider() func() *est2.State {
	return func() *est2.State {
		mem := reflect.New(p.Type.Type).Interface()
		state := &est2.State{
			Mem:       mem,
			MemPtr:    xunsafe.AsPointer(mem),
			Buffer:    est2.NewBuffer(p.bufferSize, p.escapeHTML),
			StateType: p.Type,
		}
		return state
	}
}

func (p *Planner) newExecution(root *stmt.Block) (*est2.Execution, error) {
	compute, err := p.newCompute(root)
	if err != nil {
		return nil, err
	}

	exec := est2.NewExecution(compute)
	return exec, nil
}

func (p *Planner) newCompute(root *stmt.Block) (est2.Compute, error) {
	computeNew, err := p.compileBlock(root)
	if err != nil {
		return nil, err
	}
	compute, err := computeNew(*p.Control)
	if err != nil {
		return nil, err
	}
	return compute, nil
}
