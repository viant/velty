package plan

import (
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/stmt"
)

func (p *Planner) selectorExpr(selector *expr.Select) (*op.Expression, error) {
	var err error
	expr := &op.Expression{}
	expr.Selector, err = p.SelectorExpr(selector)
	if err != nil {
		return nil, err
	}

	if expr.Selector == nil {
		id := selector.ID
		expr.Selector = est.NewSelector(id, selector.ID, nil, nil)
	}
	expr.Type = expr.Selector.Type
	return expr, nil
}

func (p *Planner) compileStmtSelector(actual *expr.Select) (est.New, error) {
	selExpr, err := p.selectorExpr(actual)
	if err != nil {
		return nil, err
	}
	return stmt.Selector(selExpr), nil
}
