package stmt

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"reflect"
	"unsafe"
)

type directAssign struct {
	x, y *op.Operand
}

func (a *directAssign) assignInt(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*int)(ret) = *(*int)(a.y.Exec(state))
	return ret
}

func (a *directAssign) assignString(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*string)(ret) = *(*string)(a.y.Exec(state))
	return ret
}

func (a *directAssign) assignFloat(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*float64)(ret) = *(*float64)(a.y.Exec(state))
	return ret
}

func (a *directAssign) assignBool(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*bool)(ret) = *(*bool)(a.y.Exec(state))
	return ret
}

func Assign(expressions ...*op.Expression) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		operands, err := op.Expressions(expressions).Operands(control)
		if err != nil {
			return nil, err
		}
		assign := &directAssign{x: operands[op.X], y: operands[op.Y]}

		switch expressions[op.Y].Type.Kind() {
		case reflect.Int:
			return assign.assignInt, nil
		case reflect.String:
			return assign.assignString, nil
		case reflect.Float64:
			return assign.assignFloat, nil
		case reflect.Bool:
			return assign.assignBool, nil
		default:
			return nil, fmt.Errorf("unsupported assign type: %s", expressions[op.Y].Type.String())
		}

	}, nil
}
