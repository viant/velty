package est

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type State struct {
	Mem       interface{}
	MemPtr    unsafe.Pointer
	StateType *Type
	Buffer    *Buffer
	Errors    []error
}

func (s *State) SetValue(k string, v interface{}) error {
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

func (s *State) IsValid() bool {
	return len(s.Errors) == 0
}

func (s *State) AddError(err error) {
	s.Errors = append(s.Errors, err)
}
