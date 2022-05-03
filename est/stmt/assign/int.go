package assign

import (
	est "github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsInt() est.Compute {
	if a.y.Comp != nil {
		return a.assignIntComp
	}

	if a.y.Sel != nil {
		return a.assignIntSelPtr
	}

	return a.assignIntLiteral
}

func (a *assign) assignIntLiteral(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*int)(ret) = *(*int)(*a.y.LiteralPtr)

	return ret
}

func (a *assign) assignIntComp(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*int)(ret) = *(*int)(a.y.Comp(state))
	return ret
}

func (a *assign) assignIntOffset(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*int)(ret) = *(*int)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignIntSelPtr(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*int)(ret) = *(*int)(a.y.Pointer(state))
	return ret
}
