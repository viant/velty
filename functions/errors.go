package functions

import "fmt"

type Errors struct {
}

func (e Errors) RegisterError(message string) (string, error) {
	return "", fmt.Errorf(message)
}

func (e Errors) AssertFloat(value interface{}) (string, error) {
	_, ok := value.(float64)
	if !ok {
		return "", fmt.Errorf("expected to got float but got %T", value)
	}

	return "", nil
}

func (e Errors) AssertInt(value interface{}) (string, error) {
	_, ok := value.(int)
	if !ok {
		return "", fmt.Errorf("expected to got int but got %T", value)
	}

	return "", nil
}

func (e Errors) AssertString(value interface{}) (string, error) {
	_, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected to got string but got %T", value)
	}

	return "", nil
}

func (e Errors) AssertBool(value interface{}) (string, error) {
	_, ok := value.(bool)
	if !ok {
		return "", fmt.Errorf("expected to got bool but got %T", value)
	}

	return "", nil
}
