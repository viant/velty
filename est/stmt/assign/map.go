package assign

import (
	"github.com/viant/velty/est"
	"reflect"
	"unsafe"
)

func (a *assign) assignAsMap() est.Compute {
	return func(state *est.State) unsafe.Pointer {
		destPtr := a.x.Exec(state)
		newMap := reflect.NewAt(a.y.Type, destPtr)
		src := a.y.ExecInterface(state)
		newMap.Elem().Set(reflect.ValueOf(src))
		return destPtr
	}
}
