package plan

import (
	"github.com/viant/velty/ast/expr"
	cexpr "github.com/viant/velty/est/expr"
	"github.com/viant/velty/est/op"
	"reflect"
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
	if resultType == nil {
		resultType = nonEmptyType(leftOperand.Type, rightOperand.Type)
	}
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

func nonEmptyType(types ...reflect.Type) reflect.Type {
	for _, r := range types {
		if r != nil {
			return r
		}
	}

	return nil
}
