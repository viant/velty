package functions

import "reflect"

var interfaceType reflect.Type

func init() {
	type foo struct {
		aField interface{}
	}

	interfaceType = reflect.ValueOf(foo{}).Field(0).Type()
}
