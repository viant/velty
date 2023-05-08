package op

import (
	"fmt"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/functions"
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

		receivers map[string]*funcReceiver
		functions map[string]*Function
		ns        map[string]interface{}
	}

	funcReceiver struct {
		rType reflect.Type
		index map[string]int
		funcs []*Func
	}

	Func struct {
		Name       string
		XType      *xunsafe.Type
		Literal    unsafe.Pointer
		ResultType reflect.Type
		Function   Funeexpression

		maxArgs    int
		isVariadic bool
		caller     reflect.Value
	}

	Function struct {
		Handler     interface{}
		ResultTyper func(call *expr.Call) (reflect.Type, error)
	}

	KindFunction interface {
		Kind() []reflect.Kind
		Handler() interface{}
	}

	ResultTyper interface {
		ResultType(receiver reflect.Type, call *expr.Call) (reflect.Type, error)
	}

	KindIndex struct {
		index            map[reflect.Kind]int
		functionsIndexes []*FunctionsIndex
	}

	FunctionsIndex struct {
		index   map[string]int
		methods []KindFunction
	}

	TypeFunc struct {
		Name       string
		Handler    interface{}
		ResultType reflect.Type
	}
)

func (r *funcReceiver) registerFunc(aFunc *Func) error {
	if _, ok := r.index[aFunc.Name]; ok {
		return fmt.Errorf("func %v already is defined on %v", aFunc.Name, r.rType.String())
	}

	r.index[aFunc.Name] = len(r.funcs)
	r.funcs = append(r.funcs, aFunc)

	return nil
}

func (r *funcReceiver) aFunc(id string) (*Func, bool) {
	i, ok := r.index[id]
	if !ok {
		return nil, false
	}

	return r.funcs[i], true
}

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

	result, err := f.execute(operands, state)
	if err != nil {
		return nil, err
	}

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

func (f *Func) execute(operands []*Operand, state *est.State) ([]reflect.Value, error) {
	if len(operands) >= f.maxArgs && !f.isVariadic {
		return nil, fmt.Errorf("too many non-variadic function arguments")
	}

	switch len(operands) {
	case 0:
		return f.caller.Call([]reflect.Value{}), nil
	case 1:
		return f.caller.Call([]reflect.Value{
			f.ensureValue(operands[0].ExecInterface(state), operands[0].Type),
		}), nil

	case 2:
		return f.caller.Call([]reflect.Value{
			f.ensureValue(operands[0].ExecInterface(state), operands[0].Type),
			f.ensureValue(operands[1].ExecInterface(state), operands[1].Type),
		}), nil

	case 3:
		return f.caller.Call([]reflect.Value{
			f.ensureValue(operands[0].ExecInterface(state), operands[0].Type),
			f.ensureValue(operands[1].ExecInterface(state), operands[1].Type),
			f.ensureValue(operands[2].ExecInterface(state), operands[2].Type),
		}), nil

	case 4:
		return f.caller.Call([]reflect.Value{
			f.ensureValue(operands[0].ExecInterface(state), operands[0].Type),
			f.ensureValue(operands[1].ExecInterface(state), operands[1].Type),
			f.ensureValue(operands[2].ExecInterface(state), operands[2].Type),
			f.ensureValue(operands[3].ExecInterface(state), operands[3].Type),
		}), nil

	default:
		values := make([]reflect.Value, 0, len(operands))
		for i := 0; i < len(operands); i++ {
			if i >= f.maxArgs && !f.isVariadic {
				return nil, fmt.Errorf("too many non-variadic function arguments")
			}

			anInterface := operands[i].ExecInterface(state)
			values = append(values, f.ensureValue(anInterface, operands[i].Type))
		}
	}
	return nil, nil
}

func (f *Func) ensureValue(anInterface interface{}, t reflect.Type) reflect.Value {
	if anInterface == nil {
		return reflect.Zero(t)
	}

	return reflect.ValueOf(anInterface)
}

func (f *Func) tryDiscoverReceiver(receiver interface{}, operands []*Operand, state *est.State, receiverValue reflect.Value) (func() (interface{}, error), bool) {
	if operands[0].LiteralPtr != nil {
		return nil, false
	}

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
					ifaces[i] = operands[i].ExecInterface(state)
				}

				return handler(ifaces...)
			}, true
		}
	}

	return nil, false
}

func NewFunctions(options ...interface{}) *Functions {
	var typeLookup functions.TypeParser
	for _, option := range options {
		switch actual := option.(type) {
		case functions.TypeParser:
			typeLookup = actual
		}
	}

	result := EmptyFunctions()
	_ = result.RegisterFuncNs(functions.FuncStrings, functions.Strings{})
	_ = result.RegisterFuncNs(functions.FuncMath, functions.Math{})
	_ = result.RegisterFuncNs(functions.FuncStrconv, functions.Strconv{})
	_ = result.RegisterFuncNs(functions.FuncSlices, functions.Slices{})
	_ = result.RegisterFuncNs(functions.FuncTypes, functions.Types{})
	_ = result.RegisterFuncNs(functions.FuncErrors, functions.Errors{})
	_ = result.RegisterFuncNs(functions.FuncTime, functions.Time{})
	_ = result.RegisterFuncNs(functions.FuncMaps, functions.Maps{})
	_ = result.RegisterFuncNs(functions.FuncJSON, functions.NewJSON(typeLookup))
	_ = result.RegisterFunctionKind(functions.MapHasKey, functions.HasKeyFunc)
	_ = result.RegisterFunctionKind(functions.SliceIndexBy, functions.SliceIndexByFunc)

	return result
}

func EmptyFunctions() *Functions {
	result := &Functions{
		index: map[string]int{},
		kindIndex: &KindIndex{
			index: map[reflect.Kind]int{},
		},
		funcs:     make([]*Func, 0),
		receivers: map[string]*funcReceiver{},
		ns:        map[string]interface{}{},
		functions: map[string]*Function{},
	}
	return result
}

func (f *Functions) RegisterFunction(name string, function interface{}) error {
	aFunc, err := f.NewFunc(name, function, nil)
	if err != nil {
		return err
	}

	return f.registerFunc(name, aFunc)
}

func (f *Functions) NewFunc(name string, function interface{}, resultType reflect.Type) (*Func, error) {
	if discoveredFn, rType, discovered := f.discover(nil, function); discovered {
		return &Func{
			Name:       name,
			Function:   discoveredFn,
			ResultType: rType,
			XType:      xunsafe.NewType(rType),
			Literal:    xunsafe.AsPointer(function),
		}, nil
	}

	var fType reflect.Type
	switch actual := function.(type) {
	case reflect.Type:
		fType = actual
	default:
		fType = reflect.TypeOf(actual)
	}

	if fType.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected func, got %v", function)
	}

	return f.reflectFunc(name, function, fType, resultType)
}

func (f *Functions) reflectFunc(name string, function interface{}, funcType reflect.Type, resultType reflect.Type) (*Func, error) {
	caller := reflect.ValueOf(function)

	if resultType == nil && funcType.NumOut() != 0 {
		resultType = funcType.Out(0)
	}

	if err := validateMethodSignature(funcType, resultType != nil); err != nil {
		return nil, err
	}

	aFunc := &Func{
		Name:       name,
		caller:     caller,
		ResultType: resultType,
		XType:      xunsafe.NewType(resultType),
		isVariadic: caller.Type().IsVariadic(),
		maxArgs:    caller.Type().NumIn() + 1, //reflect.Method.Call require to pass a receiver as first Arg.
		Literal:    xunsafe.AsPointer(function),
	}

	aFunc.Function = aFunc.callFunc
	return aFunc, nil
}

func validateMethodSignature(funcType reflect.Type, resultTypeSpecified bool) error {
	if funcType.NumOut() > 2 || funcType.NumOut() == 0 {
		return fmt.Errorf("function has to return one or two results ")
	}

	if funcType.Out(0).Kind() == reflect.Interface && !resultTypeSpecified {
		return fmt.Errorf("if method retunrns interface, the result type has to be specified with op.ResultTyper")
	}

	if funcType.NumOut() == 2 {
		if _, found := funcType.Out(1).MethodByName("Error"); !found {
			return fmt.Errorf("2nd return has to be an error if specified")
		}
	}
	return nil
}

func (f *Functions) registerFunc(name string, function *Func) error {
	if function.Function == nil {
		return fmt.Errorf("function not specified")
	}

	f.index[name] = len(f.funcs)
	f.funcs = append(f.funcs, function)

	return nil
}

func (f *Functions) IsFuncNs(ns string) bool {
	_, ok := f.ns[ns]
	return ok
}

func (f *Functions) Method(rType reflect.Type, id string, call *expr.Call) (*Func, error) {
	return f.method(rType, id, call)
}

func (f *Functions) method(rType reflect.Type, id string, call *expr.Call) (*Func, error) {
	switch rType {
	case nil:
		funcIndex, ok := f.index[id]
		if ok {
			return f.funcs[funcIndex], nil
		}

		function, ok := f.functions[id]
		if ok {
			return f.reflectFunc(id, function.Handler, reflect.TypeOf(function.Handler), nil)
		}

		return nil, fmt.Errorf("not found function %v", id)

	default:
		if method, ok := rType.MethodByName(id); ok {
			return f.asFunc(rType, id, method)
		}

		if method, err := f.functionByKind(id, rType, call); method != nil || err != nil {
			return method, err
		}

		return f.funcByName(rType, id)
	}
}

func (f *Functions) funcByName(rType reflect.Type, id string) (*Func, error) {
	index, ok := f.index[id]
	if ok {
		return f.funcs[index], nil
	}

	receiver := f.ensureReceiver(rType)
	if aFunc, ok := receiver.aFunc(id); ok {
		return aFunc, nil
	}

	return nil, fmt.Errorf("not found function %v for type %v", id, rType.String())
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

func (f *Functions) RegisterTypeFunc(receiverType reflect.Type, typeFunc *TypeFunc) error {
	return f.registerTypeFunc(receiverType, typeFunc)
}

func (f *Functions) RegisterStandaloneFunction(name string, function *Function) error {
	_, ok := f.functions[name]
	if ok {
		return fmt.Errorf("function %v already exists", name)
	}

	f.functions[name] = function
	return nil
}

func (f *Functions) registerTypeFunc(receiverType reflect.Type, typeFunc *TypeFunc) error {
	receiver := f.ensureReceiver(receiverType)
	_, ok := receiver.index[typeFunc.Name]
	if ok {
		return fmt.Errorf("function %v and receiver %v is already defined", typeFunc.Name, receiverType.String())
	}

	aFunc, err := f.NewFunc(typeFunc.Name, typeFunc.Handler, typeFunc.ResultType)
	if err != nil {
		return err
	}

	return receiver.registerFunc(aFunc)
}

func (f *Functions) ensureReceiver(receiverType reflect.Type) *funcReceiver {
	receiver, ok := f.receivers[asMapKey(receiverType)]
	if ok {
		return receiver
	}

	receiver = &funcReceiver{
		index: map[string]int{},
		funcs: make([]*Func, 0),
		rType: receiverType,
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

	return f.kindIndex.Add(methodName, funcDetails)
}

func (f *Functions) functionByKind(id string, rType reflect.Type, call *expr.Call) (*Func, error) {
	kind := rType.Kind()
	kindFunction, ok := f.kindIndex.KindFunction(kind, id)
	if !ok {
		return nil, nil
	}

	typer, ok := kindFunction.(ResultTyper)
	var resultType reflect.Type
	if ok {
		var err error
		resultType, err = typer.ResultType(rType, call)
		if err != nil {
			return nil, err
		}
	}

	handler := kindFunction.Handler()
	reflectFunc, err := f.reflectFunc(id, handler, reflect.TypeOf(handler), resultType)
	return reflectFunc, err
}

func (f *Functions) RegisterFuncNs(ns string, funcs interface{}) error {
	_, ok := f.ns[ns]
	if ok {
		return fmt.Errorf("%v already exists in Functions", ns)
	}

	f.ns[ns] = funcs
	return nil
}

func (f *Functions) FuncSelector(name string, parent *Selector) (*Selector, bool) {
	funcs, ok := f.ns[name]
	if !ok {
		return nil, false
	}

	return NewLiteralSelector(name, reflect.TypeOf(funcs), funcs, parent), true
}

func (f *Functions) TryDetectResultType(prev *Selector, methodName string, call *expr.Call) (reflect.Type, error) {
	if prev == nil {
		function, ok := f.functions[methodName]
		if ok && function.ResultTyper != nil {
			return function.ResultTyper(call)
		}

		return nil, nil
	}

	receiver, ok := f.ns[prev.ID]
	if !ok {
		return nil, nil
	}

	typer, ok := receiver.(MethodResultTyper)
	if ok {
		return typer.MethodResultType(methodName, call)
	}
	return nil, nil
}

func (i *KindIndex) Add(name string, details KindFunction) error {
	rType := reflect.TypeOf(details.Handler())
	if _, ok := details.(ResultTyper); ok {
		if err := validateMethodSignature(rType, true); err != nil {
			return err
		}
	} else {
		if err := validateMethodSignature(rType, false); err != nil {
			return err
		}
	}

	kinds := details.Kind()

	for _, kind := range kinds {
		functionsIndex := i.GetOrCreate(kind)
		functionsIndex.Add(name, details)
	}

	return nil
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
