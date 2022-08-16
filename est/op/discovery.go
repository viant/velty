package op

import (
	"github.com/viant/velty/est"
	"reflect"
)

type (
	//Discoveryable allows optimizing method calls.
	//if receiver T, receives args T1,T2 then it is possible to do type assertion like follow
	//actual, ok := aFunc.(func(receiver T, a1 T1, a2 T2))
	//TODO: reimplement it
	Discoveryable interface {
		Discover(aFunc interface{}) (func(operands []*Operand, state *est.State) (interface{}, error), reflect.Type, bool)
	}

	discoveryableMock struct{}
)

func (d discoveryableMock) Discover(aFunc interface{}) (Funeexpression, bool) {
	return nil, false
}
