package functions

import "fmt"

type Errors struct {
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
