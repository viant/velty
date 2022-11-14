package assign

import (
	est "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type assign struct {
	x, y *op2.Operand
}

func (a *assign) assignValue() est.Compute {
	return func(state *est.State) unsafe.Pointer {
		ptr := a.y.Exec(state)
		if ptr != nil {
			xunsafe.Copy(a.x.Exec(state), ptr, int(a.x.Type.Size()))
		}

		return ptr
	}
}

func Assign(expressions ...*op2.Expression) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		operands, err := op2.Expressions(expressions).Operands(control)
		if err != nil {
			return nil, err
		}

		assginer := &assign{x: operands[op2.X], y: operands[op2.Y]}
		if isIndirectOperand(operands[op2.X]) {
			return assginer.assignValue(), nil
		}

		switch expressions[op2.X].Type.Kind() {
		case reflect.Int:
			return assginer.assignAsInt(), nil
		case reflect.String:
			return assginer.assignAsString(), nil
		case reflect.Float64:
			return assginer.assignAsFloat(), nil
		case reflect.Bool:
			return assginer.assignAsBool(), nil
		case reflect.Map:
			return assginer.assignAsMap(), nil
		default:
			return assginer.assignValue(), nil
		}

	}, nil
}

func isIndirectOperand(operand *op2.Operand) bool {
	return operand.Sel != nil && operand.Sel.Indirect
}
