package velty

import (
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/stmt"
	"github.com/viant/xunsafe"
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
	if expression.Type != nil {
		expression.ValueField = xunsafe.NewField(reflect.StructField{Name: "TEMP", Offset: 0, Type: expression.Type})
	}
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
