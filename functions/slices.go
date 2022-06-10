package functions

import (
	"fmt"
	"reflect"
)

type Slices struct {
}

func (s Slices) Length(slice interface{}) int {
	return reflect.ValueOf(slice).Len()
}

func (s Slices) StringAt(slice interface{}, index int) (string, error) {
	if actual, ok := slice.([]string); ok {
		return actual[index], nil
	}
	return "", fmt.Errorf("unexpected slice type %T", slice)
}

func (s Slices) IntAt(slice interface{}, index int) (int, error) {
	if actual, ok := slice.([]int); ok {
		return actual[index], nil
	}
	return 0, fmt.Errorf("unexpected slice type %T", slice)
}

func (s Slices) BoolAt(slice interface{}, index int) (bool, error) {
	if actual, ok := slice.([]bool); ok {
		return actual[index], nil
	}
	return false, fmt.Errorf("unexpected slice type %T", slice)
}

func (s Slices) FloatAt(slice interface{}, index int) (float64, error) {
	if actual, ok := slice.([]float64); ok {
		return actual[index], nil
	}

	return 0, fmt.Errorf("unexpected slice type %T", slice)
}
