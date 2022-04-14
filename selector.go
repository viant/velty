package velty

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/stmt"
	"github.com/viant/velty/internal/ast/expr"
	"reflect"
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
	p.updateFieldOffset(actual, selExpr.Selector)

	return stmt.Selector(selExpr), nil
}

func (p *Planner) updateFieldOffset(actual *expr.Select, selector *op.Selector) {
	if selector == nil || selector.Indirect {
		return
	}

	xField, ok := p.Type.ValueAccessor(actual.ID)
	if !ok || xField.Type.Kind() != reflect.Struct {
		return
	}

	for selector.Parent != nil {
		selector = selector.Parent
	}

	selectorField := selector.Field
	selectorField.Offset += xField.Offset
}
