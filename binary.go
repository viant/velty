package velty

import (
	"github.com/viant/velty/ast/expr"
	eexpr "github.com/viant/velty/est/expr"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/types"
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

	unify, err := types.NormalizeAndUnify(x.Type, y.Type)
	if err != nil {
		return nil, err
	}

	x.Unify = unify.X
	y.Unify = unify.Y

	resultType := notNilType(types.NormalizeType(actual.Type()), unify.RType)
	acc := p.accumulator(resultType)
	resultExpr := &op.Expression{Selector: acc, Type: acc.Type}

	computeNew, err := eexpr.Binary(actual.Token, x, y, resultExpr)
	if err != nil {
		return nil, err
	}

	return &op.Expression{
		Type: resultType,
		New:  computeNew,
	}, nil
}

func notNilType(types ...reflect.Type) reflect.Type {
	for _, rType := range types {
		if rType != nil {
			return rType
		}
	}

	return nil
}
