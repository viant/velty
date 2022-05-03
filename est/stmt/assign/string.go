package assign

import (
	est "github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsString() est.Compute {
	if a.y.Comp != nil {
		return a.assignStringComp
	}

	if a.y.Sel != nil {
		return a.assignStringSelPtr
	}

	return a.assignStringLiteral
}

func (a *assign) assignStringComp(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*string)(ret) = *(*string)(a.y.Comp(state))
	return ret
}

func (a *assign) assignStringOffset(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*string)(ret) = *(*string)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignStringSelPtr(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*string)(ret) = *(*string)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignStringLiteral(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*string)(ret) = *(*string)(*a.y.LiteralPtr)
	return ret
}
