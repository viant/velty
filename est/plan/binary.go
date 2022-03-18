package plan

import (
	"github.com/viant/velty/ast/expr"
	cexpr "github.com/viant/velty/est/expr"
	"github.com/viant/velty/est/op"
	"reflect"
)

func (p *Planner) compileBinary(actual *expr.Binary) (*op.Expression, error) {
	x, err := p.compileExpr(actual.X)
	if err != nil {
		return nil, err
	}
	y, err := p.compileExpr(actual.Y)
	if err != nil {
		return nil, err
	}

	resultType := actual.Type()
	if resultType == nil {
		resultType = nonEmptyType(x.Type, y.Type)
	}
	acc := p.accumulator(resultType)
	resultExpr := &op.Expression{Selector: acc, Type: acc.Type}

	computeNew, err := cexpr.Binary(actual.Token, x, y, resultExpr)
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
