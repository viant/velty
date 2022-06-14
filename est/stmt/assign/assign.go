package assign

import (
	"fmt"
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
		assign := &assign{x: operands[op2.X], y: operands[op2.Y]}

		switch expressions[op2.Y].Type.Kind() {
		case reflect.Int:
			return assign.assignAsInt(), nil
		case reflect.String:
			return assign.assignAsString(), nil
		case reflect.Float64:
			return assign.assignAsFloat(), nil
		case reflect.Bool:
			return assign.assignAsBool(), nil
		default:
			return assign.assignValue(), nil
			return nil, fmt.Errorf("unsupported assign type: %s", expressions[op2.Y].Type.String())
		}

	}, nil
}
