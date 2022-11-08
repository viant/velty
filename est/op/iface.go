package op

import (
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type Interface struct {
	aMap     *Map
	aSlice   *Slice
	xOperand *Operand
}

func (i *Interface) Exec(xPtr unsafe.Pointer, state *est.State) unsafe.Pointer {
	asInterface := xunsafe.AsInterface(xPtr)

	actualValue := reflect.ValueOf(asInterface)
	switch actualValue.Type().Kind() {
	case reflect.Map:
		return i.aMap.Exec(xPtr, state)
	default:
		return i.aSlice.Exec(xPtr, state)
	}
}
