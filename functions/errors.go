package functions

import (
	"fmt"
	"reflect"
)

type Errors struct {
}

func (e Errors) DiscoverInterfaces(aFunc interface{}) (func(args ...interface{}) (interface{}, error), reflect.Type, bool) {
	switch actual := aFunc.(type) {
	case func(_ Errors, message string) (string, error):
		return func(args ...interface{}) (interface{}, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(args))
			}

			asMessage := args[0].(string)
			result, err := actual(e, asMessage)
			return result, err
		}, stringType, true

	case func(_ Errors, anArg interface{}) (string, error):
		return func(operands ...interface{}) (interface{}, error) {
			if len(operands) != 1 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			anArg := operands[0]
			result, err := actual(e, anArg)
			return result, err
		}, stringType, true

	case func(_ Errors, value interface{}) (bool, error):
		return func(operands ...interface{}) (interface{}, error) {
			if len(operands) != 1 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			anArg := operands[0].(bool)
			result, err := actual(e, anArg)
			return result, err
		}, boolType, true

	case func(_ Errors, value bool, message string) (bool, error):
		return func(operands ...interface{}) (interface{}, error) {
			if len(operands) != 2 {
				return nil, fmt.Errorf("unexpected arguments number, expected 1 got %v", len(operands))
			}

			valueArg := operands[0].(bool)
			messageArg := operands[1].(string)
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
