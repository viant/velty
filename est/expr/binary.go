package expr

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"reflect"
)

type binaryExpr struct {
	x *op.Operand
	y *op.Operand
	z *op.Operand
}

func Binary(token ast.Token, exprs ...*op.Expression) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		oprands, err := op.Expressions(exprs).Operands(control)
		if err != nil {
			return nil, err
		}
		binary := &binaryExpr{x: oprands[op.X], y: oprands[op.Y], z: oprands[op.Z]}

		switch exprs[0].Type.Kind() {
		case reflect.Int:
			return computeInt(token, binary)
		case reflect.Float64:
			return computeFloat(token, binary)
		case reflect.String:
			return computeString(token, binary)
		case reflect.Bool:
			return computeBool(token, binary)
		}
		return nil, fmt.Errorf("unsupported %v", exprs[0].Type.String())
	}, nil
}
