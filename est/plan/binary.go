package plan

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	cexpr "github.com/viant/velty/est/expr"
	"github.com/viant/velty/est/op"
	"reflect"
)

//Binary

func (p *Planner) compileBinary(actual *expr.Binary) (*op.Expression, error) {
	x, err := p.compileExpr(actual.X)
	if err != nil {
		return nil, err
	}
	y, err := p.compileExpr(actual.Y)
	if err != nil {
		return nil, err
	}
	zType := actual.Type()
	switch actual.Token {
	case ast.LEQ, ast.LSS, ast.GTR, ast.GTE:
		zType = reflect.TypeOf(true)
	}
	acc := p.accumulator(zType)
	z := &op.Expression{Selector: acc, Type: acc.Type}

	computeNew, err := cexpr.Binary(actual.Token, x, y, z)
	if err != nil {
		return nil, err
	}
	return &op.Expression{
		Type: zType,
		New:  computeNew,
	}, nil
}
