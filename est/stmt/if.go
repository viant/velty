package stmt

import (
	"fmt"
	est "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"github.com/viant/xunsafe"
	"reflect"
	"sync"
	"unsafe"
)

type If struct {
	ElseIf    est.Compute
	Block     est.Compute
	Condition *op2.Operand
}

// zero-check cache for struct types
type zeroMethodKind int

const (
	zeroNone zeroMethodKind = iota
	zeroIsZeroVal
	zeroIsZeroPtr
	zeroZeroVal
	zeroZeroPtr
)

var zeroCache sync.Map // map[reflect.Type]zeroMethodKind

func lookupZeroMethod(t reflect.Type) zeroMethodKind {
	if v, ok := zeroCache.Load(t); ok {
		return v.(zeroMethodKind)
	}
	kind := zeroNone
	// value receiver IsZero/Zero
	if m, ok := t.MethodByName("IsZero"); ok && m.Type.NumIn() == 1 && m.Type.NumOut() == 1 && m.Type.Out(0).Kind() == reflect.Bool {
		kind = zeroIsZeroVal
	} else if m, ok := t.MethodByName("Zero"); ok && m.Type.NumIn() == 1 && m.Type.NumOut() == 1 && m.Type.Out(0).Kind() == reflect.Bool {
		kind = zeroZeroVal
	} else {
		// pointer receiver methods
		pt := reflect.PtrTo(t)
		if m, ok := pt.MethodByName("IsZero"); ok && m.Type.NumIn() == 1 && m.Type.NumOut() == 1 && m.Type.Out(0).Kind() == reflect.Bool {
			kind = zeroIsZeroPtr
		} else if m, ok := pt.MethodByName("Zero"); ok && m.Type.NumIn() == 1 && m.Type.NumOut() == 1 && m.Type.Out(0).Kind() == reflect.Bool {
			kind = zeroZeroPtr
		}
	}
	zeroCache.Store(t, kind)
	return kind
}

func (i *If) computeWithoutElse(state *est.State) unsafe.Pointer {
	if *(*bool)(i.Condition.Exec(state)) {
		return i.Block(state)
	}
	return nil
}

func (i *If) compute(state *est.State) unsafe.Pointer {
	if *(*bool)(i.Condition.Exec(state)) {
		return i.Block(state)
	}
	return i.ElseIf(state)
}

func NewIf(condition *op2.Expression, block, elseIf est.New) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		result := &If{}
		var err error

		result.Condition, err = conditionOperand(condition, control)
		if err != nil {
			return nil, err
		}

		result.Block, err = block(control)
		if err != nil {
			return nil, err
		}

		if elseIf != nil {
			result.ElseIf, err = elseIf(control)
			if err != nil {
				return nil, err
			}
		}

		if elseIf == nil {
			return result.computeWithoutElse, nil
		}
		return result.compute, nil
	}, nil
}

func conditionOperand(condition *op2.Expression, control est.Control) (*op2.Operand, error) {
	anOperand, err := condition.Operand(control)
	if err != nil {
		return nil, err
	}

	var rType reflect.Type
	if condition.Type != nil {
		rType = condition.Type
	} else if condition.Selector != nil {
		rType = condition.Selector.Type
	}

	if condition.New != nil {
		rType = reflect.TypeOf(true)
	}

	if rType == nil {
		return nil, fmt.Errorf("couldn't determine type of the %v\n", condition.Selector.Name)
	}

	// Explicitly reject unsupported non-pointer map/struct before any coercion
	if rType.Kind() == reflect.Map || rType.Kind() == reflect.Struct {
		return nil, fmt.Errorf("unsupported comparable type %v", rType.Kind())
	}

	// If it's already boolean, pass through
	if rType.Kind() == reflect.Bool {
		return anOperand, err
	}

	newOperand := &op2.Operand{}
	// Do not apply generic unifier; we define explicit coercion semantics below.

	switch rType.Kind() {
	case reflect.Slice:
		xs := xunsafe.NewSlice(rType)
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			if anPtr == nil {
				return est.FalseValuePtr
			}
			if xs.Len(anPtr) > 0 {
				return est.TrueValuePtr
			}
			return est.FalseValuePtr
		}

	case reflect.String:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			stringPtr := (*string)(anPtr)
			if stringPtr != nil && len(*stringPtr) > 0 {
				return est.TrueValuePtr
			}
			return est.FalseValuePtr
		}

	case reflect.Int:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			intPtr := (*int)(anPtr)
			if intPtr != nil && *intPtr != 0 {
				return est.TrueValuePtr
			}
			return est.FalseValuePtr
		}

	case reflect.Float64:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			intPtr := (*float64)(anPtr)
			if intPtr != nil && *intPtr != 0 {
				return est.TrueValuePtr
			}
			return est.FalseValuePtr
		}

	case reflect.Ptr:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			if anPtr == nil {
				return est.FalseValuePtr
			}
			pointee := *(*unsafe.Pointer)(anPtr)
			if pointee == nil {
				return est.FalseValuePtr
			}
			// check zero method for *struct
			elem := rType.Elem()
			if elem.Kind() == reflect.Struct {
				switch lookupZeroMethod(elem) {
				case zeroIsZeroVal:
					v := reflect.NewAt(elem, pointee).Elem()
					res := v.MethodByName("IsZero").Call(nil)[0].Bool()
					if !res {
						return est.TrueValuePtr
					}
					return est.FalseValuePtr
				case zeroZeroVal:
					v := reflect.NewAt(elem, pointee).Elem()
					res := v.MethodByName("Zero").Call(nil)[0].Bool()
					if !res {
						return est.TrueValuePtr
					}
					return est.FalseValuePtr
				case zeroIsZeroPtr:
					v := reflect.NewAt(reflect.PtrTo(elem), unsafe.Pointer(&pointee)).Elem()
					res := v.MethodByName("IsZero").Call(nil)[0].Bool()
					if !res {
						return est.TrueValuePtr
					}
					return est.FalseValuePtr
				case zeroZeroPtr:
					v := reflect.NewAt(reflect.PtrTo(elem), unsafe.Pointer(&pointee)).Elem()
					res := v.MethodByName("Zero").Call(nil)[0].Bool()
					if !res {
						return est.TrueValuePtr
					}
					return est.FalseValuePtr
				}
			}
			return est.TrueValuePtr
		}

	case reflect.Bool:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			return anPtr
		}

	case reflect.Map:
		return nil, fmt.Errorf("unsupported comparable type %v", rType.Kind())
	case reflect.Struct:
		return nil, fmt.Errorf("unsupported comparable type %v", rType.Kind())
	default:
		return nil, fmt.Errorf("unsupported comparable type %v", rType.Kind())
	}

	return newOperand, nil
}
