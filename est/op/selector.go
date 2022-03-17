package op

import (
	"github.com/viant/velty/est"
	"unsafe"
)

func (e *Expression) newIndirectSelector() est.Compute {
	upstream, upstreamLen := est.Upstream(e.Selector)
	return func(state *est.State) unsafe.Pointer {
		ret := state.MemPtr
		for i := 0; i < upstreamLen; i++ {
			ret = upstream(i, ret)
		}

		ret = e.Selector.Pointer(ret)
		return ret
	}
}
