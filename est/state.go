package est

import (
	"fmt"
	"unsafe"
)

//TODO all privte p;easr
type State struct {
	Mem       interface{}
	MemPtr    unsafe.Pointer
	Index     map[string]int
	Selectors []*Selector
	Buffer    *Buffer
}

func (s State) Pointer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(s.MemPtr) + offset)
}

func (s State) SetValue(k string, v interface{}) error {
	idx, ok := s.Index[k]
	if !ok {
		return fmt.Errorf("undefined: %v", k)
	}
	sel := s.Selectors[idx]
	if sel.Primitive {
		sel.Set(s.MemPtr, v)
		return nil
	}
	sel.SetValue(s.MemPtr, v)
	return nil
}
