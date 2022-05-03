package assign

import (
	"github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsBool() est.Compute {
	if a.y.Comp != nil {
		return a.assignBoolComp
	}

	if a.y.Sel != nil {
		return a.assignBoolSelPtr
	}

	return a.assignBoolLiteral
}

func (a *assign) assignBoolComp(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*bool)(ret) = *(*bool)(a.y.Comp(state))
	return ret
}

func (a *assign) assignBoolOffset(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*bool)(ret) = *(*bool)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignBoolSelPtr(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*bool)(ret) = *(*bool)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignBoolLiteral(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*bool)(ret) = *(*bool)(*a.y.LiteralPtr)
	return ret
}
