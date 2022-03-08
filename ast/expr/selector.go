package expr

import (
	"reflect"
)

type Select struct {
	ID string
}

func (s Select) Type() reflect.Type {
	//TODO: Should selector have type? If no, how do we create i.e. expr.If
	return nil
}
