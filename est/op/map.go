package op

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/keys"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type Map struct {
	mapOperand   *Operand
	indexOperand *Operand
	isValueIface bool
	elemKind     reflect.Kind
}

func (m *Map) Exec(mapPtr unsafe.Pointer, state *est.State) unsafe.Pointer {
	aMap := m.mapOperand.AsInterface(mapPtr)
	if aMap == nil {
		return nil
	}

	execInterface := m.indexOperand.ExecInterface(state)
	anIndex := keys.Normalize(execInterface)
	mapValue := reflect.ValueOf(aMap)
	if mapValue.Kind() != reflect.Map {
		return nil
	}

	actualValue := mapValue.MapIndex(reflect.ValueOf(anIndex))
	if !actualValue.IsValid() {
		return nil
	}

	iface := actualValue.Interface()
	if m.isValueIface {
		return unsafe.Pointer(&iface)
	}

	switch m.elemKind {
	case reflect.Map, reflect.Slice:
		return xunsafe.AsPointer(actualValue.Pointer())
	default:
		return xunsafe.AsPointer(iface)
	}

	//pointer := unsafe.Pointer(actualValue.Pointer())
	//fmt.Println(*(*int)(pointer))
	//return pointer
}
