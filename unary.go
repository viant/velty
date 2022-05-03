package velty

import (
	aexpr "github.com/viant/velty/ast/expr"
	eexpr "github.com/viant/velty/est/expr"
	"github.com/viant/velty/est/op"
)

func (p *Planner) compileUnary(actual *aexpr.Unary) (*op.Expression, error) {
	x, err := p.compileExpr(actual.X)
	if err != nil {
		return nil, err
	}

	acc := p.accumulator(x.Type)

	resultExpr := &op.Expression{Selector: acc, Type: acc.Type}

	computeNew, err := eexpr.Unary(actual.Token, x, resultExpr)
	if err != nil {
		return nil, err
	}

	return &op.Expression{
		Type: x.Type,
		New:  computeNew,
	}, nil
}
