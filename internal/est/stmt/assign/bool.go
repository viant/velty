package assign

import (
	"github.com/viant/velty/internal/est"
	"unsafe"
)

func (a *assign) assignAsBool() est.Compute {
	if a.y.Comp != nil {
		return a.assignBoolComp
	}

	if a.y.Sel != nil {
		if a.y.Offset != nil {
			return a.assignBoolOffset
		}

		return a.assignBoolSelPtr
	}

	return a.assignBoolLiteral
}

func (a *assign) assignBoolComp(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*bool)(ret) = *(*bool)(a.y.Comp(state))
	return ret
}

func (a *assign) assignBoolOffset(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*bool)(ret) = *(*bool)(state.Pointer(*a.y.Offset))
	return ret
}

func (a *assign) assignBoolSelPtr(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*bool)(ret) = *(*bool)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignBoolLiteral(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*bool)(ret) = *(*bool)(*a.y.LiteralPtr)
	return ret
}
