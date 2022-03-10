package expr

import (
	"reflect"
)

type Select struct {
	ID   string
	Call *Call
}

func (s Select) Type() reflect.Type {
	return nil
}
