package functions

import "reflect"

type StaticKindFunc struct {
	kind       reflect.Kind
	handler    interface{}
	resultType reflect.Type
}

func (a *StaticKindFunc) Kind() reflect.Kind {
	return a.kind
}

func (a *StaticKindFunc) Handler() interface{} {
	return a.handler
}
