package stmt

import (
	"fmt"
	"github.com/viant/velty/est"
	"unsafe"
)

func Break() est.New {
	return func(control est.Control) (est.Compute, error) {
		if !control.InLoop() {
			return nil, fmt.Errorf("#break can only be used inside #for or #foreach")
		}
		return func(state *est.State) unsafe.Pointer {
			state.RequestBreak()
			return nil
		}, nil
	}
}
