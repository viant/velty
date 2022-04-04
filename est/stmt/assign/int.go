package assign

import (
	est2 "github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsInt() est2.Compute {
	if a.y.Comp != nil {
		return a.assignIntComp
	}

	if a.y.Sel != nil {
		if a.y.Offset != nil {
			return a.assignIntOffset
		}

		return a.assignIntSelPtr
	}

	return a.assignIntLiteral
}

func (a *assign) assignIntLiteral(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*int)(ret) = *(*int)(*a.y.LiteralPtr)

	return ret
}

func (a *assign) assignIntComp(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*int)(ret) = *(*int)(a.y.Comp(state))
	return ret
}

func (a *assign) assignIntOffset(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*int)(ret) = *(*int)(state.Pointer(*a.y.Offset))
	return ret
}

func (a *assign) assignIntSelPtr(state *est2.State) unsafe.Pointer {
	ret := state.Pointer(*a.x.Offset)
	*(*int)(ret) = *(*int)(a.y.Pointer(state))
	return ret
}
