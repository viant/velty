package plan

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/est"
)

type Stack struct {
	planners []*Planner
}

func NewStack(parent *Planner) *Stack {
	planners := make([]*Planner, 1)
	planners[0] = parent
	return &Stack{
		planners: planners,
	}
}

func (s *Stack) last() *Planner {
	return s.planners[len(s.planners)-1]
}

func (s *Stack) statementResolved(statement ast.Statement) {
	switch statement.(type) {
	case *stmt.Statement:
		return
	}

	s.planners = s.planners[:len(s.planners)-1]
}

func (s *Stack) newScope() *Planner {
	scope := s.last().NewScope()
	scope.selectors = s.newSelectors()
	scope.stack = s
	scope.Type = s.last().Type
	scope.cache = s.last().cache

	s.planners = append(s.planners, scope)

	return scope
}

func (s *Stack) newSelectors() *est.Selectors {
	selectors := est.NewSelectors()
	for _, planner := range s.planners {
		selectors.Merge(planner.selectors)
	}

	return selectors
}
