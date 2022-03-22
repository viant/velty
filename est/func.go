package est

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

var (
	boolType        = reflect.TypeOf(true)
	stringType      = reflect.TypeOf("")
	stringSliceType = reflect.TypeOf([]string{""})
	intType         = reflect.TypeOf(0)
	uint8Type       = reflect.TypeOf(uint8(0))
	float64Type     = reflect.TypeOf(0.0)
)

type (
	Functions struct {
		indexes map[string]int
		funcs   []*Func
	}

	Func struct {
		Caller     reflect.Value
		ResultType reflect.Type
		args       []func(pointer unsafe.Pointer) interface{}

		Function func(...unsafe.Pointer) (unsafe.Pointer, interface{})
		usePtrs  bool
	}
)

func (f *Func) CallPtrs(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
	return f.Function(pointers...)
}

func (f *Func) Call(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
	values := make([]reflect.Value, len(f.args))
	for i := 0; i < len(values); i++ {
		values[i] = reflect.ValueOf(f.args[i](pointers[i]))
	}

	result := f.Caller.Call(values)
	return xunsafe.EnsurePointer(result[0].Interface()), result[0].Interface()
}

func NewFunctions() *Functions {
	return &Functions{
		indexes: map[string]int{},
		funcs:   make([]*Func, 0),
	}
}

func (f *Functions) Register(name string, function interface{}) error {
	if discoveredFn, rType, discovered := f.discover(function); discovered {
		aFunc := &Func{
			Function:   discoveredFn,
			ResultType: rType,
		}

		return f.RegisterFunc(name, aFunc)
	}

	var fType reflect.Type
	switch actual := function.(type) {
	case reflect.Type:
		fType = actual
	default:
		fType = reflect.TypeOf(actual)
	}

	if fType.Kind() != reflect.Func {
		return fmt.Errorf("expected func, got %v", function)
	}

	caller := reflect.ValueOf(function)
	args := make([]func(pointer unsafe.Pointer) interface{}, fType.NumIn())
	for i := 0; i < len(args); i++ {
		args[i] = ptrAsInterface(fType.In(i))
	}

	var outType reflect.Type
	if fType.NumOut() != 0 {
		outType = fType.Out(0)
	}

	if fType.NumOut() > 2 || fType.NumOut() == 0 {
		return fmt.Errorf("function has to return one or two results ")
	}

	if fType.NumOut() == 2 {
		if _, found := fType.Out(1).MethodByName("Error"); !found {
			return fmt.Errorf("2nd return has to be an error if specified")
		}
	}

	f.indexes[name] = len(f.funcs)
	aFunc := &Func{
		Caller:     caller,
		ResultType: outType,
		args:       args,
	}

	aFunc.Function = aFunc.Call
	f.funcs = append(f.funcs, aFunc)

	return nil
}

func (f *Functions) RegisterFunc(name string, function *Func) error {
	if function.Function == nil {
		return fmt.Errorf("function not specified")
	}

	f.indexes[name] = len(f.funcs)
	f.funcs = append(f.funcs, function)

	return nil
}

func ptrAsInterface(rType reflect.Type) func(pointer unsafe.Pointer) interface{} {
	switch rType.Kind() {
	case reflect.String:
		return func(pointer unsafe.Pointer) interface{} {
			return *(*string)(pointer)
		}
	case reflect.Int:
		return func(pointer unsafe.Pointer) interface{} {
			return *(*int)(pointer)
		}
	case reflect.Float64:
		return func(pointer unsafe.Pointer) interface{} {
			return *(*int)(pointer)
		}
	case reflect.Bool:
		return func(pointer unsafe.Pointer) interface{} {
			return *(*bool)(pointer)
		}
	case reflect.Int64:
		return func(pointer unsafe.Pointer) interface{} {
			return *(*int64)(pointer)
		}
	case reflect.Uint8:
		return func(pointer unsafe.Pointer) interface{} {
			return *(*uint8)(pointer)
		}

	case reflect.Slice:
		switch rType.Elem().Kind() {
		case reflect.String:
			return func(pointer unsafe.Pointer) interface{} {
				return *(*[]string)(pointer)
			}
		case reflect.Int:
			return func(pointer unsafe.Pointer) interface{} {
				return *(*[]int)(pointer)
			}
		case reflect.Float64:
			return func(pointer unsafe.Pointer) interface{} {
				return *(*[]int)(pointer)
			}
		case reflect.Bool:
			return func(pointer unsafe.Pointer) interface{} {
				return *(*[]bool)(pointer)
			}
		case reflect.Int64:
			return func(pointer unsafe.Pointer) interface{} {
				return *(*[]int64)(pointer)
			}
		case reflect.Uint8:
			return func(pointer unsafe.Pointer) interface{} {
				return *(*[]uint8)(pointer)
			}
		}
	}

	xfield := xunsafe.NewField(reflect.StructField{Name: "ValueGetter", Type: rType})
	return func(pointer unsafe.Pointer) interface{} {
		return xfield.Value(pointer)
	}
}

func (f *Functions) ByName(id string) (*Func, bool) {
	index, ok := f.indexes[id]
	if !ok {
		return nil, false
	}

	return f.funcs[index], true
}

func (f *Functions) discover(function interface{}) (func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}), reflect.Type, bool) {
	switch actual := function.(type) {
	case func(s, substr string) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*string)(pointers[0]), *(*string)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue

		}, boolType, true

	case func(s1, s2 string) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]), *(*string)(pointers[1]))
			return unsafe.Pointer(&val), &val
		}, stringType, true

	case func(s string) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]))
			return unsafe.Pointer(&val), &val
		}, stringType, true

	case func(s1, s2 string) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]), *(*string)(pointers[1]))
			return unsafe.Pointer(&val), &val
		}, intType, true

	case func(s1 string) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]))
			return unsafe.Pointer(&val), &val
		}, intType, true

	case func(s string) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]))
			return unsafe.Pointer(&val), &val
		}, boolType, true

	case func(s1, s2 string, start int) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 3 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]), *(*string)(pointers[1]), *(*int)(pointers[2]))
			return unsafe.Pointer(&val), &val
		}, intType, true

	case func(s, old, new string) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 3 {
				return nil, nil
			}

			if actual(*(*string)(pointers[0]), *(*string)(pointers[1]), *(*string)(pointers[2])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(s, split string) []string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]), *(*string)(pointers[1]))
			return unsafe.Pointer(&val), &val
		}, stringSliceType, true

	case func(s1, s2 string, i int) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 3 {
				return nil, nil
			}

			if actual(*(*string)(pointers[0]), *(*string)(pointers[1]), *(*int)(pointers[2])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(s string, i int) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			val := actual(*(*string)(pointers[0]), *(*int)(pointers[1]))
			return unsafe.Pointer(&val), &val
		}, stringType, true

	case func(s string, i, end int) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 3 {
				return nil, nil
			}

			v := actual(*(*string)(pointers[0]), *(*int)(pointers[1]), *(*int)(pointers[2]))
			return unsafe.Pointer(&v), &v
		}, stringType, true

	case func(i []int, i2 int) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]int)(pointers[0]), *(*int)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []bool, i2 bool) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]bool)(pointers[0]), *(*bool)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []float64, i2 float64) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]float64)(pointers[0]), *(*float64)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []uint8, i2 uint8) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]uint8)(pointers[0]), *(*uint8)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []string, i2 string) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]string)(pointers[0]), *(*string)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []int, i2 []int) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]int)(pointers[0]), *(*[]int)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []bool, i2 []bool) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]bool)(pointers[0]), *(*[]bool)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []float64, i2 []float64) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]float64)(pointers[0]), *(*[]float64)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []uint8, i2 []uint8) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]uint8)(pointers[0]), *(*[]uint8)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []string, i2 []string) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]string)(pointers[0]), *(*[]string)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}
			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []int, i2 int) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			v := actual(*(*[]int)(pointers[0]), *(*int)(pointers[1]))
			return unsafe.Pointer(&v), &v
		}, intType, true

	case func(i []bool, i2 int) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			if actual(*(*[]bool)(pointers[0]), *(*int)(pointers[1])) {
				return TrueValuePtr, &trueValue
			}

			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []float64, i2 int) float64:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			v := actual(*(*[]float64)(pointers[0]), *(*int)(pointers[1]))
			return unsafe.Pointer(&v), &v
		}, float64Type, true

	case func(i []uint8, i2 int) uint8:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			v := actual(*(*[]uint8)(pointers[0]), *(*int)(pointers[1]))
			return unsafe.Pointer(&v), &v
		}, uint8Type, true

	case func(i []string, i2 int) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 2 {
				return nil, nil
			}

			v := actual(*(*[]string)(pointers[0]), *(*int)(pointers[1]))
			return unsafe.Pointer(&v), &v
		}, stringType, true

	case func(i []int) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			if actual(*(*[]int)(pointers[0])) {
				return TrueValuePtr, &trueValue
			}

			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []bool) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			if actual(*(*[]bool)(pointers[0])) {
				return TrueValuePtr, &trueValue
			}

			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []float64) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			if actual(*(*[]float64)(pointers[0])) {
				return TrueValuePtr, &trueValue
			}

			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []uint8) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			if actual(*(*[]uint8)(pointers[0])) {
				return TrueValuePtr, &trueValue
			}

			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []string) bool:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			if actual(*(*[]string)(pointers[0])) {
				return TrueValuePtr, &trueValue
			}

			return FalseValuePtr, &falseValue
		}, boolType, true

	case func(i []int) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*[]int)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, intType, true

	case func(i []bool) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*[]bool)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, intType, true

	case func(i []float64) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*[]float64)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, intType, true

	case func(i []string) int:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*[]string)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, intType, true

	case func(int2 int) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*int)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, stringType, true

	case func(int2 bool) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*bool)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, stringType, true

	case func(int2 float64) string:
		return func(pointers ...unsafe.Pointer) (unsafe.Pointer, interface{}) {
			if len(pointers) < 1 {
				return nil, nil
			}

			v := actual(*(*float64)(pointers[0]))

			return unsafe.Pointer(&v), &v
		}, stringType, true
	}

	return nil, nil, false
}
