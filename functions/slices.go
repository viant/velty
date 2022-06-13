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

func (s Slices) ReverseStrings(slice interface{}) ([]string, error) {
	stringsSlice, ok := slice.([]string)
	if !ok {
		return []string{}, fmt.Errorf("unexpected type, exptected []string but got %T", slice)
	}

	newSlice := make([]string, len(stringsSlice))
	for i, sValue := range stringsSlice {
		newSlice[len(newSlice)-1-i] = sValue
	}

	return newSlice, nil
}

func (s Slices) ReverseFloats(slice interface{}) ([]float64, error) {
	stringsSlice, ok := slice.([]float64)
	if !ok {
		return []float64{}, fmt.Errorf("unexpected type, exptected []float64 but got %T", slice)
	}

	newSlice := make([]float64, len(stringsSlice))
	for i, sValue := range stringsSlice {
		newSlice[len(newSlice)-1-i] = sValue
	}

	return newSlice, nil
}

func (s Slices) ReverseInts(slice interface{}) ([]int, error) {
	stringsSlice, ok := slice.([]int)
	if !ok {
		return []int{}, fmt.Errorf("unexpected type, exptected []float64 but got %T", slice)
	}

	newSlice := make([]int, len(stringsSlice))
	for i, sValue := range stringsSlice {
		newSlice[len(newSlice)-1-i] = sValue
	}

	return newSlice, nil
}
