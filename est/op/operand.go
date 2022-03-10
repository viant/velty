package op

import (
	"github.com/viant/velty/est"
	"unsafe"
)

type Operand struct {
	LiteralPtr *unsafe.Pointer
	Offset     *uintptr
	Sel        *est.Selector
	Comp       est.Compute
}

func (o *Operand) Pointer(mem *est.State) unsafe.Pointer {
	return unsafe.Pointer(uintptr(mem.MemPtr) + *o.Offset)
}

func (o *Operand) Exec(state *est.State) unsafe.Pointer {
	if o.LiteralPtr != nil {
		return *o.LiteralPtr
	}
	if o.Offset != nil {
		return o.Pointer(state)
	}
	if o.Sel != nil {
		//TODO this is not enought for pointer and accessors check igo
		return o.Sel.Pointer(state.MemPtr)
	}
	return o.Comp(state)
}
