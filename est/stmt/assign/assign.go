package assign

import (
	"fmt"
	est2 "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"reflect"
)

type assign struct {
	x, y *op2.Operand
}

func Assign(expressions ...*op2.Expression) (est2.New, error) {
	return func(control est2.Control) (est2.Compute, error) {
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
			return nil, fmt.Errorf("unsupported assign type: %s", expressions[op2.Y].Type.String())
		}

	}, nil
}
