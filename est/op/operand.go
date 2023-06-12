package op

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"github.com/viant/xunsafe/converter"
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
	NamedIFace bool
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
		if rType.Kind() == reflect.Interface && rType.NumMethod() > 0 {
			o.NamedIFace = true
		}
	} else {
		o.XType = nil
	}
}

func (o *Operand) ExecInterface(state *est.State) interface{} {
	valuePtr := o.Exec(state)
	return o.AsInterface(valuePtr)
}

func (o *Operand) ExecValue(state *est.State) reflect.Value {
	valuePtr := o.Exec(state)
	var anInterface interface{}
	switch o.XType.Kind() {
	case reflect.Interface:
		if o.NamedIFace {
			return reflect.NewAt(o.Type, valuePtr).Elem()
		}
		anInterface = xunsafe.AsInterface(valuePtr)
		return reflect.ValueOf(anInterface)
	case reflect.Func:
		anInterface = o.XType.Value(valuePtr)
	default:
		if o.LiteralPtr != nil {
			return reflect.ValueOf(o.XType.Value(valuePtr))
		}
		anInterface = o.XType.Interface(valuePtr)
	}
	return reflect.ValueOf(anInterface)
}

func (o *Operand) AsInterface(valuePtr unsafe.Pointer) interface{} {
	var anInterface interface{}
	switch o.XType.Kind() {
	case reflect.Interface:
		if o.NamedIFace {
			anInterface = reflect.NewAt(o.Type, valuePtr).Elem().Interface()
			fmt.Printf("LL %T %v\n", anInterface, anInterface)
			return anInterface
		}
		anInterface = xunsafe.AsInterface(valuePtr)
	case reflect.Func:
		anInterface = o.XType.Value(valuePtr)
	default:
		if o.LiteralPtr != nil {
			return o.XType.Value(valuePtr)
		}

		anInterface = o.XType.Interface(valuePtr)
	}

	return anInterface
}

func (o *Operand) AsValue(valuePtr unsafe.Pointer) interface{} {
	var anInterface interface{}
	switch o.XType.Kind() {
	case reflect.Interface:
		anInterface = xunsafe.AsInterface(valuePtr)
	default:
		anInterface = o.XType.Value(valuePtr)
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

func (o *Operand) SetUnifier(x converter.UnifyFn) {
	var unify func(pointer unsafe.Pointer) unsafe.Pointer
	if x != nil {
		unify = func(pointer unsafe.Pointer) unsafe.Pointer {
			ptr, _ := x(pointer)
			return ptr
		}
	}

	o.unify = unify
}
