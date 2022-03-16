package est

import (
	"unsafe"
)

func Upstream(selector *Selector) (func(index int, ptr unsafe.Pointer) unsafe.Pointer, int) {
	sel := selector.Parent
	counter := -1
	for sel != nil {
		sel = sel.Parent
		counter++
	}

	sel = selector.Parent
	parents := make([]*Selector, counter+1)
	for counter >= 0 {
		parents[counter] = sel
		sel = sel.Parent
		counter--
	}

	parentLen := len(parents)

	return func(index int, ptr unsafe.Pointer) unsafe.Pointer {
		ptr = parents[index].ValuePointer(ptr)
		return ptr
	}, parentLen
}
