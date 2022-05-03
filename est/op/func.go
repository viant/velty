package op

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/utils"
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

type Funeexpression = func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer

type (
	Functions struct {
		index map[string]int
		funcs []*Func

		receivers map[string]Receiver
	}

	Receiver struct {
		index map[string]int
		funcs []*Func
	}

	Method struct {
		ReceiverType reflect.Type
	}

	Func struct {
		caller     reflect.Value
		ResultType reflect.Type
		Function   Funeexpression
	}
)

func (f *Func) CallPtrs(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
	return f.Function(accumulator, operands, state)
}

func (f *Func) funcCall(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
	values := make([]reflect.Value, len(operands))
	for i := 0; i < len(values); i++ {
		ptr := operands[i].Exec(state)

		if operands[i].Sel != nil {
			values[i] = reflect.ValueOf(operands[i].Sel.Interface(state.MemPtr))
		} else {
			values[i] = reflect.ValueOf(asInterface(operands[i].Type, ptr))
		}
	}

	result := f.caller.Call(values)
	accumulator.SetValue(state.MemPtr, result[0].Interface())
	return accumulator.Pointer(state.MemPtr)
}

func NewFunctions() *Functions {
	return &Functions{
		index:     map[string]int{},
		funcs:     make([]*Func, 0),
		receivers: map[string]Receiver{},
	}
}

func (f *Functions) RegisterFunction(name string, function interface{}) error {
	name = utils.UpperCaseFirstLetter(name)

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

	aFunc, err := f.reflectFunc(function, fType)
	if err != nil {
		return err
	}

	f.index[name] = len(f.funcs) - 1
	f.funcs = append(f.funcs, aFunc)

	return nil
}

func (f *Functions) reflectFunc(function interface{}, fType reflect.Type) (*Func, error) {
	caller := reflect.ValueOf(function)

	var outType reflect.Type
	if fType.NumOut() != 0 {
		outType = fType.Out(0)
	}

	if fType.NumOut() > 2 || fType.NumOut() == 0 {
		return nil, fmt.Errorf("function has to return one or two results ")
	}

	if fType.NumOut() == 2 {
		if _, found := fType.Out(1).MethodByName("Error"); !found {
			return nil, fmt.Errorf("2nd return has to be an error if specified")
		}
	}

	aFunc := &Func{
		caller:     caller,
		ResultType: outType,
	}

	aFunc.Function = aFunc.funcCall
	return aFunc, nil
}

func (f *Functions) RegisterFunc(name string, function *Func) error {
	if function.Function == nil {
		return fmt.Errorf("function not specified")
	}

	f.index[name] = len(f.funcs)
	f.funcs = append(f.funcs, function)

	return nil
}

func (f *Functions) Method(rType reflect.Type, id string) (*Func, bool) {
	id = utils.UpperCaseFirstLetter(id)
	if method, ok := rType.MethodByName(id); ok {
		return f.asFunc(rType, id, method)
	}

	return f.funcByName(id)
}

func (f *Functions) funcByName(id string) (*Func, bool) {
	index, ok := f.index[id]
	if !ok {
		return nil, false
	}

	return f.funcs[index], true
}

func (f *Functions) asFunc(receiverType reflect.Type, id string, method reflect.Method) (*Func, bool) {
	f.ensureReceiver(receiverType)
	receiver, _ := f.receivers[asMapKey(receiverType)]
	index, ok := receiver.index[id]
	if ok {
		return receiver.funcs[index], true
	}

	methodSignature := method.Func.Interface()
	aFunc := &Func{}
	if funExpr, resultType, ok := f.discover(methodSignature); ok {
		aFunc.Function = funExpr
		aFunc.ResultType = resultType
	} else {
		var err error
		aFunc, err = f.reflectFunc(methodSignature, method.Type)
		if err != nil {
			return nil, false
		}
	}

	receiver.index[id] = len(receiver.funcs)
	receiver.funcs = append(receiver.funcs, aFunc)
	return aFunc, true
}

func (f *Functions) discover(function interface{}) (Funeexpression, reflect.Type, bool) {
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

func (f *Functions) RegisterTypeFunc(t reflect.Type, id string, function *Func) error {
	receiver := f.ensureReceiver(t)
	id = utils.UpperCaseFirstLetter(id)
	_, ok := receiver.index[id]
	if ok {
		return fmt.Errorf("function %v and receiver %v is already defined", id, t.String())
	}

	receiver.index[id] = len(receiver.funcs)
	receiver.funcs = append(receiver.funcs, function)

	return nil
}

func (f *Functions) ensureReceiver(receiverType reflect.Type) *Receiver {
	receiver, ok := f.receivers[asMapKey(receiverType)]
	if ok {
		return &receiver
	}

	receiver = Receiver{
		index: map[string]int{},
		funcs: make([]*Func, 0),
	}
	f.receivers[asMapKey(receiverType)] = receiver

	return &receiver
}

func asMapKey(receiverType reflect.Type) string {
	return receiverType.String()
}

//TODO: Move to the selector
func asInterface(t reflect.Type, pointer unsafe.Pointer) interface{} {
	switch t.Kind() {
	case reflect.Int:
		return *(*int)(pointer)
	case reflect.Float64:
		return *(*float64)(pointer)
	case reflect.Bool:
		return *(*bool)(pointer)
	case reflect.String:
		return *(*string)(pointer)
	}

	return xunsafe.AsInterface
}
