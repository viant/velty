package velty

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/parser"
	"unsafe"
)

type evaluator struct {
	x       *op.Operand
	cache   *cache
	control est.Control
	parent  *Planner
}

func (e *evaluator) compute(state *est.State) unsafe.Pointer {
	varValue := *(*string)(e.x.Exec(state))
	if cacheValue, ok := e.cache.expression(varValue); ok {
		newState := e.newState(cacheValue.planner, state)
		return cacheValue.compute(newState)
	}

	block, err := parser.Parse([]byte(varValue))
	if err != nil {
		return est.EmptyStringPtr
	}

	evaluatorPlanner := e.parent.New()
	exec, err := evaluatorPlanner.newCompute(block)
	if err != nil {
		return est.EmptyStringPtr
	}

	newState := e.newState(evaluatorPlanner, state)
	e.cache.put(varValue, evaluatorPlanner, exec)

	return exec(newState)
}

func (e *evaluator) newState(planner *Planner, state *est.State) *est.State {
	newState := planner.stateProvider()()

	for _, valueAccessor := range e.parent.Type.ValueAccessors() {
		valueAccessor.SetValue(newState.MemPtr, valueAccessor.Value(state.MemPtr))
	}

	newState.Buffer = state.Buffer
	return newState
}

func evaluate(expr *op.Expression, cache *cache, parent *Planner) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		x, err := expr.Operand(control, false)
		if err != nil {
			return nil, err
		}

		return (&evaluator{
			x:       x,
			cache:   cache,
			control: control,
			parent:  parent,
		}).compute, nil
	}, nil
}
