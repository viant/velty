package assign

import (
	est2 "github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsFloat() est2.Compute {
	if a.y.Comp != nil {
		return a.assignFloatComp
	}

	if a.y.Sel != nil {
		if a.y.Offset != nil {
			return a.assignFloatOffset
		}

		return a.assignFloatSelPtr
	}

	return a.assignFloatLiteral
}

func (a *assign) assignFloatComp(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*float64)(ret) = *(*float64)(a.y.Comp(state))
	return ret
}

func (a *assign) assignFloatOffset(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*float64)(ret) = *(*float64)(state.Pointer(*a.y.Offset))
	return ret
}

func (a *assign) assignFloatSelPtr(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*float64)(ret) = *(*float64)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignFloatLiteral(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*float64)(ret) = *(*float64)(*a.y.LiteralPtr)
	return ret
}
