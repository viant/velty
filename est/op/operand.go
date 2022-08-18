package op

import (
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type Operand struct {
	LiteralPtr *unsafe.Pointer
	Sel        *Selector
	Comp       est.Compute
	Type       reflect.Type
	XType      *xunsafe.Type
}

func (o *Operand) Pointer(state *est.State) unsafe.Pointer {
	return unsafe.Pointer(uintptr(state.MemPtr) + o.Sel.Offset + o.Sel.ParentOffset)
}

func (o *Operand) Exec(state *est.State) unsafe.Pointer {
	if o.Comp != nil {
		return o.Comp(state)
	}

	if o.LiteralPtr != nil {
		return *o.LiteralPtr
	}

	if o.Sel != nil {
		return o.Pointer(state)
	}

	return unsafe.Pointer(uintptr(state.MemPtr) + o.Sel.Offset + o.Sel.ParentOffset)
}

func (o *Operand) IsIndirect() bool {
	return o.Sel == nil || o.Sel.Indirect
}

func (o *Operand) SetType(rType reflect.Type) {
	o.Type = rType
	if rType != nil {
		o.XType = xunsafe.NewType(rType)
	} else {
		o.XType = nil
	}
}
