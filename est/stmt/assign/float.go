package assign

import (
	"github.com/viant/velty/est"
	"unsafe"
)

func (a *assign) assignAsFloat() est.Compute {
	if a.y.Comp != nil {
		return a.assignFloatComp
	}

	if a.y.Sel != nil {
		return a.assignFloatSelPtr
	}

	return a.assignFloatLiteral
}

func (a *assign) assignFloatComp(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*float64)(ret) = *(*float64)(a.y.Comp(state))
	return ret
}

func (a *assign) assignFloatOffset(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*float64)(ret) = *(*float64)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignFloatSelPtr(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*float64)(ret) = *(*float64)(a.y.Pointer(state))
	return ret
}

func (a *assign) assignFloatLiteral(state *est.State) unsafe.Pointer {
	ret := a.x.Pointer(state)
	*(*float64)(ret) = *(*float64)(*a.y.LiteralPtr)
	return ret
}
