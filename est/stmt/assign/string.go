package assign

import (
	"github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsString() est.Compute {
	if a.y.Comp != nil {
		return a.assignStringComp
	}

	if a.y.Sel != nil {
		if a.y.Offset != nil {
			return a.assignStringOffset
		}

		return a.assignStringSelPtr
	}

	return a.assignStringLiteral
}

func (a *assign) assignStringComp(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*string)(ret) = *(*string)(a.y.Comp(state))
	return ret
}

func (a *assign) assignStringOffset(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*string)(ret) = *(*string)(state.Pointer(*a.y.Offset))
	return ret
}

func (a *assign) assignStringSelPtr(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*string)(ret) = *(*string)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignStringLiteral(state *est.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*string)(ret) = *(*string)(*a.y.LiteralPtr)
	return ret
}
