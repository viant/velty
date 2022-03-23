package op

import (
	"fmt"
	"github.com/viant/velty/est"
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

type FuncExpression = func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer

type (
	Functions struct {
		indexes map[string]int
		funcs   []*Func
	}

	Func struct {
		Caller     reflect.Value
		ResultType reflect.Type

		Function FuncExpression
	}
)

func (f *Func) CallPtrs(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
	return f.Function(accumulator, operands, state)
}

func (f *Func) Call(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
	values := make([]reflect.Value, len(operands))
	for i := 0; i < len(values); i++ {
		values[i] = reflect.ValueOf(operands[i].Exec(state))
	}

	result := f.Caller.Call(values)
	accumulator.SetValue(state.MemPtr, result[0].Interface())
	return accumulator.Pointer(state.MemPtr)
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

func (f *Functions) ByName(id string) (*Func, bool) {
	index, ok := f.indexes[id]
	if !ok {
		return nil, false
	}

	return f.funcs[index], true
}

func (f *Functions) discover(function interface{}) (FuncExpression, reflect.Type, bool) {
	switch actual := function.(type) {
	case func(s, substr string) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(s1, s2 string) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, stringType, true

	case func(s string) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*string)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, stringType, true

	case func(s1, s2 string) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetInt(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, intType, true

	case func(s1 string) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetInt(state.MemPtr, actual(*(*string)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, intType, true

	case func(s string) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*string)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, boolType, true

	case func(s1, s2 string, start int) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 3 {
				return nil
			}

			accumulator.SetInt(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state)), *(*int)(operands[2].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, intType, true

	case func(s, old, new string) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 3 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state)), *(*string)(operands[2].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, boolType, true

	case func(s, split string) []string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetValue(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, stringSliceType, true

	case func(s1, s2 string, i int) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 3 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state)), *(*int)(operands[2].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, boolType, true

	case func(s string, i int) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, stringType, true

	case func(s string, i, end int) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 3 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*string)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state)), *(*int)(operands[2].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, stringType, true

	case func(i []int, i2 int) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]int)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []bool, i2 bool) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]bool)(operands[0].Exec(state)), *(*bool)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []float64, i2 float64) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]float64)(operands[0].Exec(state)), *(*float64)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []uint8, i2 uint8) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]uint8)(operands[0].Exec(state)), *(*uint8)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []string, i2 string) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []int, i2 []int) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]int)(operands[0].Exec(state)), *(*[]int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []bool, i2 []bool) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]bool)(operands[0].Exec(state)), *(*[]bool)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []float64, i2 []float64) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]float64)(operands[0].Exec(state)), *(*[]float64)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []uint8, i2 []uint8) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]uint8)(operands[0].Exec(state)), *(*[]uint8)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []string, i2 []string) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]string)(operands[0].Exec(state)), *(*[]string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []int, i2 int) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetInt(state.MemPtr, actual(*(*[]int)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, intType, true

	case func(i []bool, i2 int) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]bool)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []float64, i2 int) float64:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetFloat64(state.MemPtr, actual(*(*[]float64)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, float64Type, true

	case func(i []uint8, i2 int) uint8:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetUint8(state.MemPtr, actual(*(*[]uint8)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, uint8Type, true

	case func(i []string, i2 int) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*[]string)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, stringType, true

	case func(i []int) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]int)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []bool) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]bool)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []float64) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]float64)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []uint8) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]uint8)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []string) bool:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.Set(state.MemPtr, actual(*(*[]string)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, boolType, true

	case func(i []int) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetInt(state.MemPtr, actual(*(*[]int)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, intType, true

	case func(i []bool) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetInt(state.MemPtr, actual(*(*[]bool)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, intType, true

	case func(i []float64) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetValue(state.MemPtr, actual(*(*[]float64)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, intType, true

	case func(i []string) int:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetValue(state.MemPtr, actual(*(*[]string)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, intType, true

	case func(int2 int) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*int)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		}, stringType, true

	case func(int2 bool) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*bool)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, stringType, true

	case func(int2 float64) string:
		return func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 1 {
				return nil
			}

			accumulator.SetString(state.MemPtr, actual(*(*float64)(operands[0].Exec(state))))
			return accumulator.Pointer(state.MemPtr)

		}, stringType, true
	}

	return nil, nil, false
}
