package est

import (
	"fmt"
	"unsafe"
)

type State struct {
	Mem       interface{}
	MemPtr    unsafe.Pointer
	index     map[string]int
	selectors []*Selector
	Buffer    *Buffer
}

func (s State) SetValue(k string, v interface{}) error {
	idx, ok := s.index[k]
	if !ok {
		return fmt.Errorf("undefined: %v", k)
	}
	sel := s.selectors[idx]
	if sel.Primitive {
		sel.Set(s.MemPtr, v)
		return nil
	}
	sel.SetValue(s.MemPtr, v)
	return nil
}
