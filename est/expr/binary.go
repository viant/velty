package expr

import (
	"fmt"
	est2 "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"github.com/viant/velty/internal/ast"
	"reflect"
)

type binaryExpr struct {
	x *op2.Operand
	y *op2.Operand
	z *op2.Operand
}

func Binary(token ast.Token, exprs ...*op2.Expression) (est2.New, error) {
	return func(control est2.Control) (est2.Compute, error) {
		oprands, err := op2.Expressions(exprs).Operands(control)
		if err != nil {
			return nil, err
		}

		binary := &binaryExpr{x: oprands[op2.X], y: oprands[op2.Y], z: oprands[op2.Z]}
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
