package stmt

import (
	"fmt"
	est "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"github.com/viant/xunsafe/converter"
	"reflect"
	"unsafe"
)

type If struct {
	ElseIf    est.Compute
	Block     est.Compute
	Condition *op2.Operand
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

	unify, err := converter.Unify(reflect.TypeOf(true), rType)
	if err != nil {
		return nil, err
	}

	if err != nil || rType.Kind() == reflect.Bool {
		return anOperand, err
	}

	newOperand := &op2.Operand{}
	if unify.Y != nil {
		newOperand.SetUnifier(unify.Y)
		rType = reflect.TypeOf(true)
	}

	switch rType.Kind() {
	case reflect.Slice:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			anHeader := (*reflect.SliceHeader)(anPtr)
			if anHeader != nil && anHeader.Len > 0 {
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
			if anPtr != nil {
				return est.TrueValuePtr
			}
			return est.FalseValuePtr
		}

	case reflect.Bool:
		newOperand.Comp = func(state *est.State) unsafe.Pointer {
			anPtr := anOperand.Exec(state)
			return anPtr
		}

	default:
		return nil, fmt.Errorf("unsupported comparable type %v", condition.Type.Kind())
	}

	return newOperand, nil
}
