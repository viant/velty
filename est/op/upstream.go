package op

import (
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

func Upstream(selector *Selector, derefLast bool) func(state *est.State) unsafe.Pointer {
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

	callers := make([]func(accumulator *Selector, selectors []*Operand, state *est.State) unsafe.Pointer, parentLen)
	for i := 0; i < parentLen; i++ {
		if parents[i].Func == nil {
			continue
		}
		callers[i] = parents[i].Func.CallFunc
	}

	return func(state *est.State) unsafe.Pointer {
		ptr := state.MemPtr
		if ptr == nil {
			return zeroValuePtr
		}

		for i := 0; i < parentLen; i++ {
			if parents[i].Func != nil {
				ptr = callers[i](parents[i], parents[i].Args, state)
			} else if parents[i].Slice != nil {
				ptr = parents[i].Slice.Exec(state)
			} else {
				if (derefLast && i == parentLen-1) || (i < parentLen-1 && callers[i+1] != nil) {
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
