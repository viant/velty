package expr

import (
	"reflect"
)

type Range struct {
	X *Literal
	Y *Literal
}

func (r *Range) Type() reflect.Type {
	if r.X.Type() != nil {
		return r.X.Type()
	}
	return r.Y.Type()
}
