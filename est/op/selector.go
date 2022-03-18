package op

import (
	"github.com/viant/velty/est"
	"unsafe"
)

func (e *Expression) newIndirectSelector() est.Compute {
	upstream := est.Upstream(e.Selector)
	return func(state *est.State) unsafe.Pointer {
		ret := upstream(state.MemPtr)
		return ret
	}
}
