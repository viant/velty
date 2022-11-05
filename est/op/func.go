package op

import (
	"fmt"
	"github.com/viant/velty/est"
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
		index         map[string]int
		kindIndex     *KindIndex
		kindFunctions []KindFunction
		funcs         []*Func

		receivers map[string]*Receiver
	}

	Receiver struct {
		index map[string]int
		funcs []*Func
	}

	Func struct {
		caller     reflect.Value
		ResultType reflect.Type
		Function   Funeexpression

		maxArgs    int
		isVariadic bool
		Name       string
		XType      *xunsafe.Type
	}

	KindFunction interface {
		Kind() reflect.Kind
		Handler() interface{}
	}

	ResultTyper interface {
		ResultType(receiver reflect.Type) (reflect.Type, error)
	}

	KindIndex struct {
		index            map[reflect.Kind]int
		functionsIndexes []*FunctionsIndex
	}

	FunctionsIndex struct {
		index   map[string]int
		methods []KindFunction
	}
)

func (f *Func) CallFunc(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
	anIface, err := f.Function(operands, state)
	if err != nil {
		state.AddError(err)
	}

	if anIface != nil {
		accumulator.SetValue(state.MemPtr, anIface)
		if f.XType.Type().Kind() == reflect.Map {
			return unsafe.Pointer(reflect.ValueOf(f.XType.Ref(anIface)).Pointer())
		}

		return xunsafe.AsPointer(anIface)
	}

	return nil
}

func (f *Func) callFunc(operands []*Operand, state *est.State) (interface{}, error) {
	if len(operands) == 0 {
		return nil, fmt.Errorf("expected to got min 1 operand but got %v", len(operands))
	}

	receiverIface := operands[0].AsInterface(state)
	receiverValue := reflect.ValueOf(receiverIface)
	if handler, ok := f.tryDiscoverReceiver(receiverIface, operands, state, receiverValue); ok {
		return handler()
	}

	values := make([]reflect.Value, len(operands))
	values[0] = receiverValue

	for i := 1; i < len(values); i++ {
		if i >= f.maxArgs && !f.isVariadic {
			return nil, fmt.Errorf("too many non-variadic function arguments")
		}

		anInterface := operands[i].AsInterface(state)
		if anInterface == nil {
			values[i] = reflect.Zero(operands[i].Type)
		} else {
			values[i] = reflect.ValueOf(anInterface)
		}
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

func (f *Func) tryDiscoverReceiver(receiver interface{}, operands []*Operand, state *est.State, receiverValue reflect.Value) (func() (interface{}, error), bool) {
	if actual, ok := receiver.(Discoveryable); ok {
		method := receiverValue.MethodByName(f.Name)
		handler, _, ok := actual.Discover(method.Interface())
		if ok {
			return func() (interface{}, error) {
				return handler(operands, state)
			}, true
		}
	}

	if actual, ok := receiver.(DiscoveryableIface); ok {
		method := receiverValue.MethodByName(f.Name)
		handler, _, ok := actual.DiscoverInterfaces(method.Interface())
		if ok {
			return func() (interface{}, error) {
				ifaces := make([]interface{}, len(operands))
				ifaces[0] = receiver
				for i := 1; i < len(operands); i++ {
					ifaces[i] = operands[i].AsInterface(state)
				}

				return handler(ifaces...)
			}, true
		}
	}

	return nil, false
}

func NewFunctions() *Functions {
	return &Functions{
		index: map[string]int{},
		kindIndex: &KindIndex{
			index: map[reflect.Kind]int{},
		},
		funcs:     make([]*Func, 0),
		receivers: map[string]*Receiver{},
	}
}

func (f *Functions) RegisterFunction(name string, function interface{}) error {
	if discoveredFn, rType, discovered := f.discover(nil, function); discovered {
		aFunc := &Func{
			Name:       name,
			Function:   discoveredFn,
			ResultType: rType,
			XType:      xunsafe.NewType(rType),
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

	aFunc, err := f.reflectFunc(name, function, fType, nil)
	if err != nil {
		return err
	}

	f.index[name] = len(f.funcs) - 1
	f.funcs = append(f.funcs, aFunc)

	return nil
}

func (f *Functions) reflectFunc(name string, function interface{}, funcType reflect.Type, resultType reflect.Type) (*Func, error) {
	caller := reflect.ValueOf(function)

	if resultType == nil && funcType.NumOut() != 0 {
		resultType = funcType.Out(0)
	}

	if funcType.NumOut() > 2 || funcType.NumOut() == 0 {
		return nil, fmt.Errorf("function has to return one or two results ")
	}

	if funcType.NumOut() == 2 {
		if _, found := funcType.Out(1).MethodByName("Error"); !found {
			return nil, fmt.Errorf("2nd return has to be an error if specified")
		}
	}

	aFunc := &Func{
		Name:       name,
		caller:     caller,
		ResultType: resultType,
		XType:      xunsafe.NewType(resultType),
		isVariadic: caller.Type().IsVariadic(),
		maxArgs:    caller.Type().NumIn() + 1, //reflect.Method.Call require to pass a receiver as first Arg.
	}

	aFunc.Function = aFunc.callFunc
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

func (f *Functions) Method(rType reflect.Type, id string) (*Func, error) {
	if method, ok := rType.MethodByName(id); ok {
		return f.asFunc(rType, id, method)
	}

	if method, err := f.functionByKind(id, rType); method != nil || err != nil {
		return method, err
	}

	return f.funcByName(id)
}

func (f *Functions) funcByName(id string) (*Func, error) {
	index, ok := f.index[id]
	if !ok {
		return nil, fmt.Errorf("no such function %v", id)
	}

	return f.funcs[index], nil
}

func (f *Functions) asFunc(receiverType reflect.Type, id string, method reflect.Method) (*Func, error) {
	f.ensureReceiver(receiverType)
	receiver, _ := f.receivers[asMapKey(receiverType)]
	index, ok := receiver.index[id]
	if ok {
		return receiver.funcs[index], nil
	}

	methodSignature := method.Func.Interface()
	aFunc := &Func{}
	if funExpr, resultType, ok := f.discover(receiverType, methodSignature); ok {
		aFunc.Function = funExpr
		aFunc.ResultType = resultType
	} else {
		var err error
		aFunc, err = f.reflectFunc(id, methodSignature, method.Type, nil)
		if err != nil {
			return nil, err
		}
	}

	receiver.index[id] = len(receiver.funcs)
	receiver.funcs = append(receiver.funcs, aFunc)
	return aFunc, nil
}

func (f *Functions) discover(receiverType reflect.Type, function interface{}) (Funeexpression, reflect.Type, bool) {
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

func incorrectArgumentsError(wanted string, got []*Operand) error {
	return fmt.Errorf("expected to got %v but got %v", wanted, len(got))
}

func (f *Functions) RegisterFunctionKind(methodName string, funcDetails KindFunction) error {
	handler := funcDetails.Handler()
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		return fmt.Errorf("unexpected function handler type, expected Function, got %T", handler)
	}

	f.kindIndex.Add(methodName, funcDetails)
	return nil
}

func (f *Functions) functionByKind(id string, rType reflect.Type) (*Func, error) {
	kind := rType.Kind()
	kindFunction, ok := f.kindIndex.KindFunction(kind, id)
	if !ok {
		return nil, nil
	}

	typer, ok := kindFunction.(ResultTyper)
	var resultType reflect.Type
	if ok {
		var err error
		resultType, err = typer.ResultType(rType)
		if err != nil {
			return nil, err
		}
	}

	handler := kindFunction.Handler()
	reflectFunc, err := f.reflectFunc(id, handler, reflect.TypeOf(handler), resultType)
	return reflectFunc, err
}

func (i *KindIndex) Add(name string, details KindFunction) {
	functionsIndex := i.GetOrCreate(details.Kind())
	functionsIndex.Add(name, details)
}

func (i *KindIndex) GetOrCreate(kind reflect.Kind) *FunctionsIndex {
	functionsIndex, done := i.Get(kind)
	if done {
		return functionsIndex
	}

	result := &FunctionsIndex{index: map[string]int{}}
	i.index[kind] = len(i.functionsIndexes)
	i.functionsIndexes = append(i.functionsIndexes, result)

	return result
}

func (i *KindIndex) Get(kind reflect.Kind) (*FunctionsIndex, bool) {
	index, ok := i.index[kind]
	if ok {
		return i.functionsIndexes[index], true
	}
	return nil, false
}

func (i *KindIndex) KindFunction(kind reflect.Kind, id string) (KindFunction, bool) {
	functionsIndex, ok := i.Get(kind)
	if !ok {
		return nil, false
	}

	return functionsIndex.Get(id)
}

func (i *FunctionsIndex) Add(name string, details KindFunction) {
	i.index[name] = len(i.methods)
	i.methods = append(i.methods, details)
}

func (i *FunctionsIndex) Get(id string) (KindFunction, bool) {
	index, ok := i.index[id]
	if !ok {
		return nil, false
	}

	return i.methods[index], true
}
