package op

import (
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"github.com/viant/xunsafe/converter"
	"reflect"
	"unsafe"
)

func Upstream(selector *Selector, derefLast bool, refLast bool) func(state *est.State) unsafe.Pointer {
	derefLast = derefLast || converter.IsPrimitive(selector.Type)
	sel := selector.Parent
	counter := -1
	for sel != nil {
		sel = sel.Parent
		counter++
	}

	sel = selector.Parent
	parents := make([]*Selector, counter+2)
	for counter >= 0 {
		parents[counter] = sel
		sel = sel.Parent
		counter--
	}

	parents[len(parents)-1] = selector
	parentLen := len(parents)

	var zeroValuePtr unsafe.Pointer
	var value interface{}

	switch selector.Type.Kind() {
	case reflect.Bool:
		zeroValuePtr = est.FalseValuePtr
	case reflect.String:
		zeroValuePtr = est.EmptyStringPtr
	case reflect.Int:
		zeroValuePtr = est.ZeroIntPtr
	case reflect.Float64:
		zeroValuePtr = est.ZeroFloatPtr
	default:
		value = reflect.New(selector.Type).Interface()
		zeroValuePtr = xunsafe.AsPointer(value)
	}

	shouldRefLast := selector.Type.Kind() == reflect.Ptr

	return func(state *est.State) unsafe.Pointer {
		ptr := state.MemPtr
		if ptr == nil {
			return zeroValuePtr
		}

		for i := 0; i < parentLen; i++ {
			shouldRef := shouldRefLast && (i == parentLen-1)
			if parents[i].Literal != nil {
				ptr = refIfNeeded(parents[i].Literal, shouldRef)
			} else if parents[i].Func != nil {
				args := parents[i].Args
				if i != 0 { //receiver call
					args = make([]*Operand, len(parents[i].Args)) // have to copy args and replace first operand,
					// because CallFunc would call Upstream once again,
					// calling all methods multiple times
					copy(args, parents[i].Args)
					newArg := *args[0]
					newArg.Comp = nil
					if parents[i-1].Func != nil {
						newArg.Value = newArg.AsValue(ptr)
					} else {
						newArg.Value = newArg.AsInterface(ptr)

					}

					args[0] = &newArg
				}

				ptr = refIfNeeded(parents[i].Func.CallFunc(parents[i], args, state), shouldRef)
			} else if parents[i].Slice != nil {
				ptr = refIfNeeded(parents[i].Slice.Exec(ptr, state), shouldRef)
			} else if parents[i].Map != nil {
				ptr = refIfNeeded(parents[i].Map.Exec(ptr, state), shouldRef)
			} else if parents[i].InterfaceExec != nil {
				ptr = refIfNeeded(parents[i].InterfaceExec.Exec(ptr, state), shouldRef)
			} else {
				if ((!derefLast || shouldRef) && i == parentLen-1) || (i < parentLen-1 && parents[i+1].Func != nil) {
					ptr = parents[i].Pointer(ptr)
				} else {
					ptr = parents[i].ValuePointer(ptr)
				}
			}

			if ptr == nil {
				return zeroValuePtr
			}
		}

		return ptr
	}
}

func refIfNeeded(literal unsafe.Pointer, ref bool) unsafe.Pointer {
	if ref {
		literal = xunsafe.RefPointer(literal)
	}

	return literal
}
