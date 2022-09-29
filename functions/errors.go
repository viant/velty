package functions

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"reflect"
)

type Errors struct {
}

func (e Errors) Discover(aFunc interface{}) (func(operands []*op.Operand, state *est.State) (interface{}, error), reflect.Type, bool) {
	switch actual := aFunc.(type) {
	case func(_ Errors, message string) (string, error):
		return func(operands []*op.Operand, state *est.State) (interface{}, error) {
			if len(operands) != 1 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			asMessage := *(*string)(operands[0].Exec(state))
			result, err := actual(e, asMessage)
			return result, err
		}, stringType, true

	case func(_ Errors, anArg interface{}) (string, error):
		return func(operands []*op.Operand, state *est.State) (interface{}, error) {
			if len(operands) != 1 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			anArg := op.AsInterface(operands[0], operands[0].Exec(state))
			result, err := actual(e, anArg)
			return result, err
		}, stringType, true

	case func(_ Errors, value interface{}) (bool, error):
		return func(operands []*op.Operand, state *est.State) (interface{}, error) {
			if len(operands) != 1 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			anArg := *(*bool)(operands[0].Exec(state))
			result, err := actual(e, anArg)
			return result, err
		}, boolType, true

	case func(_ Errors, value bool, message string) (bool, error):
		return func(operands []*op.Operand, state *est.State) (interface{}, error) {
			if len(operands) != 2 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			valueArg := *(*bool)(operands[0].Exec(state))
			messageArg := *(*string)(operands[1].Exec(state))
			result, err := actual(e, valueArg, messageArg)
			return result, err
		}, boolType, true

	}

	return nil, nil, false
}

func (e Errors) Raise(message string) (string, error) {
	return "", fmt.Errorf(message)
}

func (e Errors) RegisterError(message string) (string, error) {
	return "", fmt.Errorf(message)
}

func (e Errors) AssertFloat(value interface{}) (bool, error) {
	_, ok := value.(float64)
	if !ok {
		return false, fmt.Errorf("expected to got float but got %T", value)
	}

	return true, nil
}

func (e Errors) AssertInt(value interface{}) (bool, error) {
	_, ok := value.(int)
	if !ok {
		return false, fmt.Errorf("expected to got int but got %T", value)
	}

	return true, nil
}

func (e Errors) AssertString(value interface{}) (bool, error) {
	_, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("expected to got string but got %T", value)
	}

	return true, nil
}

func (e Errors) AssertBool(value interface{}) (bool, error) {
	_, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("expected to got bool but got %T", value)
	}

	return true, nil
}

func (e Errors) AssertWithMessage(isValid bool, message string) (bool, error) {
	if isValid {
		return true, nil
	}

	return false, fmt.Errorf(message)
}
