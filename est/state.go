package est

import (
	"fmt"
	"reflect"
	"unsafe"
)

//TODO all privte p;easr
type State struct {
	Mem       interface{}
	MemPtr    unsafe.Pointer
	Selectors *Selectors
	Buffer    *Buffer
}

func (s *State) Pointer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(s.MemPtr) + offset)
}

func (s *State) SetValue(k string, v interface{}) error {
	idx, ok := s.Selectors.Index[k]
	if !ok {
		return fmt.Errorf("undefined: %v", k)
	}

	sel := s.Selectors.Selector(idx)
	if !sel.Indirect && sel.Kind() != reflect.Struct {
		sel.Set(s.MemPtr, v)
		return nil
	}

	sel.SetValue(s.MemPtr, v)
	return nil
}

func (s *State) Reset() {
	s.Buffer.Reset()
}
