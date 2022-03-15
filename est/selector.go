package est

import (
	"github.com/viant/xunsafe"
	"reflect"
)

//Selector represent data selector
type Selector struct {
	ID string
	*xunsafe.Field
	Primitive bool
	Parent    *Selector
}

//NewSelector create a selector
func NewSelector(id, name string, sType reflect.Type, parent *Selector) *Selector {
	return &Selector{
		ID:     id,
		Field:  &xunsafe.Field{Name: name, Type: sType},
		Parent: parent,
	}
}
