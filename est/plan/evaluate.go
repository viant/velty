package plan

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/plan/scope"
	"github.com/viant/velty/parser"
	"unsafe"
)

type evaluator struct {
	x       *op.Operand
	cache   *Cache
	control est.Control

	parent *Planner
}

func (e *evaluator) compute(state *est.State) unsafe.Pointer {
	varValue := *(*string)(e.x.Exec(state))
	if cacheValue, ok := e.cache.Expression(varValue); ok {
		newState, err := e.newState(cacheValue.Planner, state)
		if err != nil {
			return est.EmptyStringPtr
		}
		return cacheValue.Compute(newState)
	}

	block, err := parser.Parse([]byte(varValue))
	if err != nil {
		return est.EmptyStringPtr
	}

	evaluatorPlanner, err := e.evaluatorPlanner(state)

	if err != nil {
		return est.EmptyStringPtr
	}

	exec, err := evaluatorPlanner.newCompute(block)
	if err != nil {
		return est.EmptyStringPtr
	}

	newState, err := e.newState(evaluatorPlanner, state)
	if err != nil {
		return est.EmptyStringPtr
	}

	e.cache.Put(varValue, evaluatorPlanner, exec)

	return exec(newState)
}

func (e *evaluator) evaluatorPlanner(state *est.State) (*Planner, error) {
	evaluatorScope := New(0)
	scopeType := scope.NewType()
	evaluatorScope.Type = scopeType

	var err error
	for _, selector := range e.parent.selectors.Selectors() {
		if selector.Parent != nil {
			continue
		}

		if err = evaluatorScope.DefineVariable(selector.Name, selector.Value(state.MemPtr)); err != nil {
			return nil, err
		}
	}

	return evaluatorScope, nil
}

func (e *evaluator) newState(planner *Planner, state *est.State) (*est.State, error) {
	newState := planner.stateProvider()()

	var err error
	for _, selector := range e.parent.selectors.Selectors() {
		if selector.Parent != nil {
			continue
		}

		if err = newState.SetValue(selector.ID, selector.Value(state.MemPtr)); err != nil {
			return nil, err
		}
	}

	newState.Buffer = state.Buffer
	return newState, nil
}

func Evaluate(expr *op.Expression, cache *Cache, parent *Planner) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		x, err := expr.Operand(control)
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
