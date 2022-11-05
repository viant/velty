package op

import (
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type Map struct {
	mapOperand   *Operand
	indexOperand *Operand
	isValueIface bool
}

func (m *Map) Exec(state *est.State) unsafe.Pointer {
	aMap := m.mapOperand.AsInterface(state)
	if aMap == nil {
		return nil
	}

	anIndex := m.indexOperand.AsInterface(state)
	mapValue := reflect.ValueOf(aMap)
	if mapValue.Kind() != reflect.Map {
		return nil
	}

	actualValue := mapValue.MapIndex(reflect.ValueOf(anIndex))
	iface := actualValue.Interface()
	if m.isValueIface {
		return unsafe.Pointer(&iface)
	}

	return xunsafe.AsPointer(iface)
}
