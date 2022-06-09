package functions

import (
	"reflect"
)

type Slices struct {
}

func (s Slices) Length(slice interface{}) int {
	return reflect.ValueOf(slice).Len()
}
