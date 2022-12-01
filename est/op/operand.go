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
	unify      func(pointer unsafe.Pointer) unsafe.Pointer
	Repeat     bool
}

func (o *Operand) Pointer(state *est.State) unsafe.Pointer {
	return o.unifyPtr(o.pointer(state))
}

func (o *Operand) unifyPtr(pointer unsafe.Pointer) unsafe.Pointer {
	if pointer != nil && o.unify != nil {
		pointer = o.unify(pointer)
	}
	return pointer
}

func (o *Operand) pointer(state *est.State) unsafe.Pointer {
	if o.LiteralPtr != nil {
		return *o.LiteralPtr
	}

	return unsafe.Pointer(uintptr(state.MemPtr) + o.Sel.Offset + o.Sel.ParentOffset)
}

func (o *Operand) Exec(state *est.State) unsafe.Pointer {
	return o.unifyPtr(o.exec(state))
}

func (o *Operand) exec(state *est.State) unsafe.Pointer {
	if o.Comp != nil {
		return o.Comp(state)
	}

	if o.LiteralPtr != nil {
		return *o.LiteralPtr
	}

	return o.pointer(state)
}

func (o *Operand) IsIndirect() bool {
	return (o.Sel == nil && o.LiteralPtr == nil) || (o.Sel != nil && o.Sel.Indirect)
}

func (o *Operand) SetType(rType reflect.Type) {
	o.Type = rType
	if rType != nil {
		o.XType = xunsafe.NewType(o.getXType(rType))
	} else {
		o.XType = nil
	}
}

func (o *Operand) ExecInterface(state *est.State) interface{} {
	valuePtr := o.Exec(state)
	return o.AsInterface(valuePtr)
}

func (o *Operand) AsInterface(valuePtr unsafe.Pointer) interface{} {
	var anInterface interface{}
	switch o.XType.Kind() {
	case reflect.Interface:
		anInterface = xunsafe.AsInterface(valuePtr)
	default:
		anInterface = o.XType.Interface(valuePtr)
	}

	return anInterface
}

func (o *Operand) trySetType(rType reflect.Type) {
	if o.Type == nil {
		o.SetType(rType)
	}
}

func (o *Operand) getXType(rType reflect.Type) reflect.Type {
	return rType
}
