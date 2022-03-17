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
	Index     map[string]int
	Selectors []*Selector
	Buffer    *Buffer
}

func (s *State) Pointer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(s.MemPtr) + offset)
}

func (s *State) SetValue(k string, v interface{}) error {
	idx, ok := s.Index[k]
	if !ok {
		return fmt.Errorf("undefined: %v", k)
	}

	sel := s.Selectors[idx]
	if !sel.Indirect && sel.Kind() != reflect.Struct {
		sel.Set(s.MemPtr, v)
		return nil
	}

	sel.SetValue(s.MemPtr, v)
	return nil
}

func (s *State) Reset() {
	for _, sel := range s.Selectors {
		switch sel.Type.Kind() {
		case reflect.Int:
			sel.SetInt(s.MemPtr, 0)
		case reflect.Int64:
			sel.SetInt64(s.MemPtr, 0)
		case reflect.Uint64:
			sel.SetUint64(s.MemPtr, 0)

		case reflect.String:
			sel.SetString(s.MemPtr, "")
		case reflect.Float64:
			sel.SetFloat64(s.MemPtr, 0.0)
		case reflect.Bool:
			sel.SetBool(s.MemPtr, false)
		}
	}

	s.Buffer.Reset()
}
