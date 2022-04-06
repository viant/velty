package expr

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/internal/ast"
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
		indirect := binary.x.IsIndirect() || binary.y.IsIndirect()

		switch exprs[0].Type.Kind() {
		case reflect.Int:
			return computeInt(token, binary, indirect)
		case reflect.Float64:
			return computeFloat(token, binary, indirect)
		case reflect.String:
			return computeString(token, binary, indirect)
		case reflect.Bool:
			return computeBool(token, binary, indirect)
		}
		return nil, fmt.Errorf("unsupported %v", exprs[0].Type.String())
	}, nil
}
