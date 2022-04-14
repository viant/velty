package expr

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/internal/ast"
	"reflect"
)

type unaryExpr struct {
	x *op.Operand
	y *op.Operand
}

func Unary(token ast.Token, exprs ...*op.Expression) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		oprands, err := op.Expressions(exprs).Operands(control)
		if err != nil {
			return nil, err
		}

		unary := &unaryExpr{x: oprands[op.X], y: oprands[op.Y]}
		indirect := unary.x.IsIndirect() || unary.y.IsIndirect()

		switch exprs[0].Type.Kind() {
		case reflect.Bool:
			return computeUnaryBool(token, unary, indirect)
		}
		return nil, fmt.Errorf("unsupported %v as unary expression", exprs[0].Type.String())
	}, nil
}
