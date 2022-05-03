package velty

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/stmt"
	"github.com/viant/velty/internal/ast/expr"
)

func (p *Planner) selectorExpr(selector *expr.Select) (*op.Expression, error) {
	var err error
	expression := &op.Expression{}
	expression.Selector, err = p.selector(selector)
	if err != nil {
		return nil, err
	}

	if expression.Selector == nil {
		id := selector.ID
		expression.Selector = op.NewSelector(id, selector.ID, nil, nil)
		expression.Selector.Placeholder = selector.FullName
	}
	expression.Type = expression.Selector.Type
	return expression, nil
}

func (p *Planner) compileStmtSelector(actual *expr.Select) (est.New, error) {
	selExpr, err := p.selectorExpr(actual)
	if err != nil {
		return nil, err
	}

	p.Type.ValueAccessor(actual.ID)
	return stmt.Selector(selExpr), nil
}
