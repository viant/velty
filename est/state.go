package est

import (
	"fmt"
	"github.com/viant/velty/utils"
	"reflect"
	"strings"
	"unsafe"
)

//TODO all privte p;easr
type State struct {
	Mem       interface{}
	MemPtr    unsafe.Pointer
	StateType *Type
	Buffer    *Buffer
	Errors    []error
}

//func (s *State) Pointer(offset uintptr) unsafe.Pointer {
//	return unsafe.Pointer(uintptr(s.MemPtr) + offset)
//}

func (s *State) SetValue(k string, v interface{}) error {
	k = utils.UpperCaseFirstLetter(k)

	xField, ok := s.StateType.ValueAccessor(k)
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

func (s *State) EmbedValue(v interface{}) error {
	k := strings.Split(reflect.TypeOf(v).String(), ".")[1]
	k = utils.UpperCaseFirstLetter(k)

	xField, ok := s.StateType.ValueAccessor(k)
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
	s.Errors = nil
}
