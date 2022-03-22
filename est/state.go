package est

import (
	"fmt"
	"github.com/viant/velty/est/plan/scope"
	"reflect"
	"unsafe"
)

//TODO all privte p;easr
type State struct {
	Mem       interface{}
	MemPtr    unsafe.Pointer
	StateType *scope.Type
	Buffer    *Buffer
}

func (s *State) Pointer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(s.MemPtr) + offset)
}

func (s *State) SetValue(k string, v interface{}) error {
	xField, ok := s.StateType.Mutator(k)
	if !ok {
		return fmt.Errorf("undefined: %v", k)
	}

	switch xField.Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		xField.SetValue(s.MemPtr, v)
	default:
		xField.Set(s.MemPtr, v)

	}

	return nil
}

func (s *State) Reset() {
	s.Buffer.Reset()
}
