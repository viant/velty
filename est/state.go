package est

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

type TemplateError error
type State struct {
	sync.Mutex
	Mem          interface{}
	MemPtr       unsafe.Pointer
	StateType    *Type
	Buffer       *Buffer
	Errors       []error
	PanicOnError bool
	isTaken      bool
}

func (s *State) SetValue(k string, v interface{}) error {
	xField, ok := s.StateType.ValueAccessor(k)
	if !ok {
		return fmt.Errorf("undefined: %v", k)
	}

	switch xField.Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Slice, reflect.Map, reflect.Interface:
		xField.SetValue(s.MemPtr, v)
	default:
		xField.Set(s.MemPtr, v)
	}

	return nil
}

func (s *State) EmbedValue(v interface{}) error {
	vType := reflect.TypeOf(v)
	if vType.Kind() == reflect.Ptr && vType.Elem().Name() != "" {
		vType = vType.Elem()
	}
	var holderName string
	if vType.Name() == "" {
		var ok bool
		holderName, ok = s.StateType.AnonymousHolder(vType)
		if !ok {
			return fmt.Errorf("not found holder for %T", v)
		}
	} else {
		holderName = strings.Split(vType.String(), ".")[1]
	}

	xField, ok := s.StateType.ValueAccessor(holderName)
	if !ok {
		return fmt.Errorf("undefined: %v", holderName)
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
	s.Lock()
	defer s.Unlock()

	s.Buffer.Reset()
	s.Errors = nil
	s.isTaken = true
}

func (s *State) IsValid() bool {
	return len(s.Errors) == 0
}

func (s *State) AddError(err error) {
	err = TemplateError(err)

	s.Errors = append(s.Errors, err)
	if s.PanicOnError {
		panic(err)
	}
}

func (s *State) Take() bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if s.isTaken {
		return false
	}

	s.isTaken = true
	return true
}
