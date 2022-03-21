package assign

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"reflect"
)

type assign struct {
	x, y *op.Operand
}

func Assign(expressions ...*op.Expression) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		operands, err := op.Expressions(expressions).Operands(control)
		if err != nil {
			return nil, err
		}
		assign := &assign{x: operands[op.X], y: operands[op.Y]}

		switch expressions[op.Y].Type.Kind() {
		case reflect.Int:
			return assign.assignAsInt(), nil
		case reflect.String:
			return assign.assignAsString(), nil
		case reflect.Float64:
			return assign.assignAsFloat(), nil
		case reflect.Bool:
			return assign.assignAsBool(), nil
		default:
			return nil, fmt.Errorf("unsupported assign type: %s", expressions[op.Y].Type.String())
		}

	}, nil
}
