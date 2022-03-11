package plan

import (
	"github.com/viant/velty/ast/expr"
	cexpr "github.com/viant/velty/est/expr"
	"github.com/viant/velty/est/op"
)

//Binary

func (p *Planner) compileBinary(actual *expr.Binary) (*op.Expression, error) {
	leftOperand, err := p.compileExpr(actual.X)
	if err != nil {
		return nil, err
	}
	rightOperand, err := p.compileExpr(actual.Y)
	if err != nil {
		return nil, err
	}

	resultType := actual.Type()
	acc := p.accumulator(resultType)
	resultExpr := &op.Expression{Selector: acc, Type: acc.Type}

	computeNew, err := cexpr.Binary(actual.Token, leftOperand, rightOperand, resultExpr)
	if err != nil {
		return nil, err
	}

	return &op.Expression{
		Type: resultType,
		New:  computeNew,
	}, nil
}
