package est

import (
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

func Upstream(selector *Selector) func(ptr unsafe.Pointer) unsafe.Pointer {
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

	switch selector.Kind() {
	case reflect.Bool:
		zeroValuePtr = FalseValuePtr
	case reflect.String:
		zeroValuePtr = EmptyStringPtr
	case reflect.Int:
		zeroValuePtr = ZeroIntPtr
	case reflect.Float64:
		zeroValuePtr = ZeroFloatPtr
	default:
		value = reflect.New(selector.Type).Interface()
		zeroValuePtr = xunsafe.AsPointer(value)
	}

	return func(ptr unsafe.Pointer) unsafe.Pointer {
		if ptr == nil {
			return zeroValuePtr
		}

		ret := ptr
		for i := 0; i < parentLen; i++ {
			ret = parents[i].ValuePointer(ret)
			if ret == nil {
				return zeroValuePtr
			}
		}

		return ret
	}
}
