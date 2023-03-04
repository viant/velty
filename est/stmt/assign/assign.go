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
	a.x.Sel.Field.MustBeAssignable(a.y.Type)

	rType := a.x.Type
	wasPtr := false
	for rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		wasPtr = true
	}

	switch rType.Kind() {
	case reflect.Int, reflect.Uint64, reflect.Int64:

		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)

				*(*int)(destPtr) = *(*int)(srcPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)

				if srcPtr != nil {
					*(**int)(destPtr) = *(**int)(srcPtr)
				}

				return srcPtr
			}
		}

	case reflect.Int8, reflect.Uint8:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*int8)(destPtr) = *(*int8)(srcPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				if srcPtr != nil {
					*(**int8)(destPtr) = *(**int8)(srcPtr)
				}

				return srcPtr
			}
		}
	case reflect.Int16, reflect.Uint16:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*int16)(destPtr) = *(*int16)(srcPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				if srcPtr != nil {
					*(**int16)(destPtr) = *(**int16)(srcPtr)
				}

				return srcPtr
			}
		}
	case reflect.Int32, reflect.Uint32:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*int32)(destPtr) = *(*int32)(srcPtr)
				return srcPtr
			}
		}

		return func(state *est.State) unsafe.Pointer {
			destPtr := a.x.Exec(state)
			srcPtr := a.y.Exec(state)
			if srcPtr != nil {
				*(**int32)(destPtr) = *(**int32)(srcPtr)
			}
			return srcPtr
		}
	case reflect.String:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*string)(destPtr) = *(*string)(srcPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				if srcPtr != nil {
					*(**string)(destPtr) = *(**string)(srcPtr)
				}
				return srcPtr
			}
		}
	case reflect.Float64:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*float64)(destPtr) = *(*float64)(srcPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				if srcPtr != nil {
					*(**float64)(destPtr) = *(**float64)(srcPtr)
				}
				return srcPtr
			}
		}
	case reflect.Float32:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*float32)(destPtr) = *(*float32)(destPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)

				if srcPtr != nil {
					*(**float32)(destPtr) = *(**float32)(destPtr)
				}
				return srcPtr
			}
		}
	case reflect.Bool:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				*(*bool)(destPtr) = *(*bool)(destPtr)
				return srcPtr
			}
		} else {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				srcPtr := a.y.Exec(state)
				if srcPtr != nil {
					*(**bool)(destPtr) = *(**bool)(destPtr)
				}
				return srcPtr
			}
		}

	case reflect.Struct:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				ptr := a.y.Exec(state)
				if ptr != nil {
					xunsafe.Copy(a.x.Exec(state), ptr, int(a.x.Type.Size()))
				}
				return ptr
			}
		}
	case reflect.Slice:
		if !wasPtr {
			return func(state *est.State) unsafe.Pointer {
				destPtr := a.x.Exec(state)
				ptr := a.y.Exec(state)

				if destPtr != nil && ptr != nil {
					sourceHeader := (*reflect.SliceHeader)(ptr)
					destHader := (*reflect.SliceHeader)(destPtr)
					destHader.Data = sourceHeader.Data
					destHader.Len = sourceHeader.Len
					destHader.Cap = sourceHeader.Cap
				}

				return ptr
			}
		}

	case reflect.Func:
		return func(state *est.State) unsafe.Pointer {
			dest := a.x.Exec(state)
			src := a.y.Exec(state)

			if dest == nil || src == nil {
				return nil
			}

			a.x.Sel.Field.SetFuncAt(dest, a.y.AsInterface(src))
			return dest
		}
	}

	if wasPtr {
		return func(state *est.State) unsafe.Pointer {
			dest := a.x.Exec(state)
			src := a.y.Exec(state)

			if src != nil && dest != nil {
				*(*unsafe.Pointer)(dest) = src
			}

			return src
		}
	}

	return func(state *est.State) unsafe.Pointer {
		dest := a.x.Exec(state)
		src := a.y.Exec(state)

		if src != nil {
			xunsafe.Copy(dest, src, int(rType.Size()))
		}

		return src
	}
}

func Assign(expressions ...*op2.Expression) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		operands, err := op2.Expressions(expressions).Operands(control, false)
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
