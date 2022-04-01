package op

import (
	est2 "github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

func Upstream(selector *Selector) func(state *est2.State) unsafe.Pointer {
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
		zeroValuePtr = est2.FalseValuePtr
	case reflect.String:
		zeroValuePtr = est2.EmptyStringPtr
	case reflect.Int:
		zeroValuePtr = est2.ZeroIntPtr
	case reflect.Float64:
		zeroValuePtr = est2.ZeroFloatPtr
	default:
		value = reflect.New(selector.Type).Interface()
		zeroValuePtr = xunsafe.AsPointer(value)
	}

	callers := make([]func(accumulator *Selector, selectors []*Operand, state *est2.State) unsafe.Pointer, parentLen)
	for i := 0; i < parentLen; i++ {
		if parents[i].Func == nil {
			continue
		}
		callers[i] = parents[i].Func.Function
	}

	return func(state *est2.State) unsafe.Pointer {
		ptr := state.MemPtr
		if ptr == nil {
			return zeroValuePtr
		}

		for i := 0; i < parentLen; i++ {
			if parents[i].Func == nil {
				ptr = parents[i].ValuePointer(ptr)
			} else {
				ptr = callers[i](parents[i], parents[i].Args, state)
			}

			if ptr == nil {
				return zeroValuePtr
			}
		}

		return ptr
	}
}
