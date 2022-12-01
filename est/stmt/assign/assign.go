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
	rType := a.x.Type
	for rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}

	switch rType.Kind() {
	case reflect.Int, reflect.Uint64, reflect.Int64:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)

			*(*int)(destPtr) = *(*int)(srcPtr)
			return srcPtr
		}
	case reflect.Int8, reflect.Uint8:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*int8)(destPtr) = *(*int8)(srcPtr)
			return srcPtr
		}
	case reflect.Int16, reflect.Uint16:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*int16)(destPtr) = *(*int16)(srcPtr)
			return srcPtr
		}
	case reflect.Int32, reflect.Uint32:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*int32)(destPtr) = *(*int32)(srcPtr)
			return srcPtr
		}
	case reflect.String:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*string)(destPtr) = *(*string)(srcPtr)
			return srcPtr
		}
	case reflect.Float64:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*float64)(destPtr) = *(*float64)(srcPtr)
			return srcPtr
		}
	case reflect.Float32:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*float32)(destPtr) = *(*float32)(destPtr)
			return srcPtr
		}
	case reflect.Bool:
		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			*(*bool)(destPtr) = *(*bool)(destPtr)
			return srcPtr
		}

	case reflect.Struct, reflect.Slice:
		return func(state *est.State) unsafe.Pointer {
			ptr := a.y.Exec(state)
			if ptr != nil {
				xunsafe.Copy(a.x.Exec(state), ptr, int(a.x.Type.Size()))
			}
			return ptr
		}

	default:
		return func(state *est.State) unsafe.Pointer {
			a.x.Sel.Field.SetValue(a.x.Exec(state), a.y.ExecInterface(state))
			return a.y.Exec(state)
		}
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
