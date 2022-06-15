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

type Funeexpression = func(operands []*Operand, state *est.State) (interface{}, error)

type (
	Functions struct {
		index map[string]int
		funcs []*Func

		receivers map[string]*Receiver
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

		maxArgs    int
		isVariadic bool
	}
)

func (f *Func) CallFunc(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
	anIface, err := f.Function(operands, state)
	if err != nil {
		state.Errors = append(state.Errors, err)
	}

	if anIface != nil {
		accumulator.SetValue(state.MemPtr, anIface)
		return xunsafe.AsPointer(anIface)
	}

	return nil
}

func (f *Func) funcCall(operands []*Operand, state *est.State) (interface{}, error) {
	values := make([]reflect.Value, len(operands))
	var argSelector *Selector

	for i := 0; i < len(values); i++ {
		valuePtr := operands[i].Exec(state)
		if i < f.maxArgs {
			argSelector = operands[i].Sel
		}

		if i >= f.maxArgs && !f.isVariadic {
			return nil, fmt.Errorf("too many non-variadic function arguments")
		}

		var anInterface interface{}
		if argSelector != nil && argSelector.ValueField != nil {
			anInterface = argSelector.ValueField.Interface(valuePtr)
		} else {
			anInterface = asInterface(operands[i].Type, valuePtr)
		}

		if anInterface == nil {
			anInterface = interface{}(nil)
		}

		values[i] = reflect.ValueOf(anInterface)
	}

	result := f.caller.Call(values)
	if len(result) == 1 {
		return result[0].Interface(), nil
	} else if len(result) == 2 {
		var err error
		var ok bool
		errInface := result[1].Interface()
		if errInface != nil {
			err, ok = errInface.(error)
			if !ok {
				return nil, fmt.Errorf("unexpected error type %T", errInface)
			}
		}

		return result[0].Interface(), err
	} else if len(result) == 0 {
		return nil, nil
	} else {
		return nil, fmt.Errorf("unexpected number of returned values, expected <= 2, but got %v", len(result))
	}
}

func NewFunctions() *Functions {
	return &Functions{
		index:     map[string]int{},
		funcs:     make([]*Func, 0),
		receivers: map[string]*Receiver{},
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
		isVariadic: caller.Type().IsVariadic(),
		maxArgs:    caller.Type().NumIn(),
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
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("(string, string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(s1, s2 string) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("(string, string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))), nil
		}, stringType, true

	case func(s string) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("(string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state))), nil
		}, stringType, true

	case func(s1, s2 string) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("(string, string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))), nil
		}, intType, true

	case func(s1 string) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("(string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state))), nil
		}, intType, true

	case func(s string) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("(string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state))), nil
		}, boolType, true

	case func(s1, s2 string, start int) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 3 {
				return nil, incorrectArgumentsError("(string, string, int)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state)), *(*int)(operands[2].Exec(state))), nil
		}, intType, true

	case func(s, old, new string) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 3 {
				return nil, incorrectArgumentsError("(string, string, string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state)), *(*string)(operands[2].Exec(state))), nil
		}, boolType, true

	case func(s, split string) []string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("(string, string)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))), nil
		}, stringSliceType, true

	case func(s1, s2 string, i int) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 3 {
				return nil, incorrectArgumentsError("(string, string, int)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state)), *(*int)(operands[2].Exec(state))), nil
		}, boolType, true

	case func(s string, i int) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("(string, int)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, stringType, true

	case func(s string, i, end int) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 3 {
				return nil, incorrectArgumentsError("(string, int, int)", operands)
			}

			return actual(*(*string)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state)), *(*int)(operands[2].Exec(state))), nil

		}, stringType, true

	case func(i []int, i2 int) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]int, []int)", operands)
			}

			return actual(*(*[]int)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(i []bool, i2 bool) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]bool, []bool)", operands)
			}

			return actual(*(*[]bool)(operands[0].Exec(state)), *(*bool)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(i []float64, i2 float64) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]float64, []float64)", operands)
			}

			return actual(*(*[]float64)(operands[0].Exec(state)), *(*float64)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(i []uint8, i2 uint8) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]uint8, []uint8)", operands)
			}

			return actual(*(*[]uint8)(operands[0].Exec(state)), *(*uint8)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(i []string, i2 string) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]string, []string)", operands)
			}

			return actual(*(*[]string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(i []int, i2 []int) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]int, []int)", operands)
			}

			return actual(*(*[]int)(operands[0].Exec(state)), *(*[]int)(operands[1].Exec(state))), nil

		}, boolType, true

	case func([]bool, []bool) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]bool, []bool)", operands)
			}

			return actual(*(*[]bool)(operands[0].Exec(state)), *(*[]bool)(operands[1].Exec(state))), nil

		}, boolType, true

	case func([]float64, []float64) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]float64, []float64)", operands)
			}

			return actual(*(*[]float64)(operands[0].Exec(state)), *(*[]float64)(operands[1].Exec(state))), nil

		}, boolType, true

	case func([]uint8, []uint8) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]uint8, []uint8)", operands)
			}

			return actual(*(*[]uint8)(operands[0].Exec(state)), *(*[]uint8)(operands[1].Exec(state))), nil

		}, boolType, true

	case func(i []string, i2 []string) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]string, []string)", operands)
			}

			return actual(*(*[]string)(operands[0].Exec(state)), *(*[]string)(operands[1].Exec(state))), nil

		}, boolType, true

	case func([]int, int) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]int, int)", operands)
			}

			return actual(*(*[]int)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, intType, true

	case func([]bool, int) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]bool, int)", operands)
			}

			return actual(*(*[]bool)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, boolType, true

	case func([]float64, int) float64:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]float64, int)", operands)
			}

			return actual(*(*[]float64)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, float64Type, true

	case func([]uint8, int) uint8:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]uint8, int)", operands)
			}

			return actual(*(*[]uint8)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, uint8Type, true

	case func([]string, int) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 2 {
				return nil, incorrectArgumentsError("([]string, int)", operands)
			}

			return actual(*(*[]string)(operands[0].Exec(state)), *(*int)(operands[1].Exec(state))), nil

		}, stringType, true

	case func([]int) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]int)", operands)
			}

			return actual(*(*[]int)(operands[0].Exec(state))), nil

		}, boolType, true

	case func([]bool) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]bool)", operands)
			}

			return actual(*(*[]bool)(operands[0].Exec(state))), nil

		}, boolType, true

	case func([]float64) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]float64)", operands)
			}

			return actual(*(*[]float64)(operands[0].Exec(state))), nil

		}, boolType, true

	case func([]uint8) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]uint8)", operands)
			}

			return actual(*(*[]uint8)(operands[0].Exec(state))), nil

		}, boolType, true

	case func([]string) bool:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]string)", operands)
			}

			return actual(*(*[]string)(operands[0].Exec(state))), nil

		}, boolType, true

	case func([]int) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]int)", operands)
			}

			return actual(*(*[]int)(operands[0].Exec(state))), nil

		}, intType, true

	case func([]bool) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]bool)", operands)
			}

			return actual(*(*[]bool)(operands[0].Exec(state))), nil

		}, intType, true

	case func([]float64) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]float64)", operands)
			}

			return actual(*(*[]float64)(operands[0].Exec(state))), nil

		}, intType, true

	case func([]string) int:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("([]string)", operands)
			}

			return actual(*(*[]string)(operands[0].Exec(state))), nil
		}, intType, true

	case func(int) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("(int)", operands)
			}

			return actual(*(*int)(operands[0].Exec(state))), nil
		}, stringType, true

	case func(bool) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("(bool)", operands)
			}

			return actual(*(*bool)(operands[0].Exec(state))), nil

		}, stringType, true

	case func(float64) string:
		return func(operands []*Operand, state *est.State) (interface{}, error) {
			if len(operands) < 1 {
				return nil, incorrectArgumentsError("(float)", operands)
			}

			return actual(*(*float64)(operands[0].Exec(state))), nil

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
		return receiver
	}

	receiver = &Receiver{
		index: map[string]int{},
		funcs: make([]*Func, 0),
	}

	f.receivers[asMapKey(receiverType)] = receiver
	return receiver
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

	return xunsafe.AsInterface(pointer)
}

func incorrectArgumentsError(wanted string, got []*Operand) error {
	return fmt.Errorf("expected to got %v but got %v", wanted, len(got))
}
