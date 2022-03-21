package op

import (
	"github.com/viant/velty/est"
	"reflect"
	"unsafe"
)

type Operand struct {
	LiteralPtr *unsafe.Pointer
	Offset     *uintptr
	Sel        *est.Selector
	Comp       est.Compute
	Type       reflect.Type
}

func (o *Operand) Pointer(state *est.State) unsafe.Pointer {
	return unsafe.Pointer(uintptr(state.MemPtr) + *o.Offset)
}

func (o *Operand) Exec(state *est.State) unsafe.Pointer {
	if o.Comp != nil {
		return o.Comp(state)
	}

	if o.LiteralPtr != nil {
		return *o.LiteralPtr
	}

	if o.Offset != nil {
		return o.Pointer(state)
	}

	return o.Sel.Pointer(state.MemPtr)
}

func (o *Operand) IsIndirect() bool {
	return (o.Sel != nil && o.Sel.Indirect) || o.Offset == nil
}
