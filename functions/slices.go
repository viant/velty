package functions

import (
	"fmt"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/keys"
	"github.com/viant/xreflect"
	"github.com/viant/xunsafe"
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

var SliceIndexByFunc = &indexSliceByFunc{}

type indexSliceByFunc struct{}

func (in *indexSliceByFunc) Kind() []reflect.Kind {
	return []reflect.Kind{reflect.Slice}
}

func (in *indexSliceByFunc) Handler() interface{} {
	return func(slice interface{}, field string) (interface{}, error) {
		sliceType := reflect.TypeOf(slice)
		if sliceType.Kind() != reflect.Slice {
			return nil, fmt.Errorf("unsupported IndexBy receiver, got %T", slice)
		}

		elemType := sliceType.Elem()
		upstream := in.upstream(elemType)
		xField, err := in.Field(upstream, elemType, field)
		if err != nil {
			return nil, err
		}

		mapType, err := in.ResultType(sliceType, nil)
		if err != nil {
			return nil, err
		}

		resultMap := reflect.MakeMap(mapType)
		if slice == nil {
			return resultMap.Interface(), nil
		}

		xSlice := xunsafe.NewSlice(sliceType)
		slicePtr := xunsafe.AsPointer(slice)

		sliceLen := xSlice.Len(slicePtr)
		for i := 0; i < sliceLen; i++ {
			sliceValueAt := xSlice.ValueAt(slicePtr, i)
			fieldValue := sliceValueAt
			for _, upstreamType := range upstream {
				fieldValue = upstreamType.Deref(fieldValue)
			}

			key := keys.Normalize(fieldValue)
			fieldValue = xField.Value(xunsafe.AsPointer(key))
			resultMap.SetMapIndex(reflect.ValueOf(fieldValue), reflect.ValueOf(sliceValueAt))
		}
		resultMapIface := resultMap.Interface()
		if err != nil {
			return nil, err
		}

		return resultMapIface, nil
	}
}

func (in *indexSliceByFunc) ResultType(receiver reflect.Type, _ *expr.Call) (reflect.Type, error) {
	if receiver.Kind() != reflect.Slice {
		return nil, fmt.Errorf("unsupported IndexBy receiver type %s", receiver.String())
	}

	return reflect.MapOf(xreflect.InterfaceType, receiver.Elem()), nil
}

func (in *indexSliceByFunc) upstream(sliceType reflect.Type) []*xunsafe.Type {
	var types []*xunsafe.Type
	for sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
		types = append(types, xunsafe.NewType(sliceType))
	}

	return types
}

func (in *indexSliceByFunc) Field(upstream []*xunsafe.Type, elemType reflect.Type, field string) (*xunsafe.Field, error) {
	var xField *xunsafe.Field
	if len(upstream) > 0 {
		xField = xunsafe.FieldByName(upstream[len(upstream)-1].Type(), field)
	} else {
		xField = xunsafe.FieldByName(elemType, field)
	}

	if xField != nil {
		return xField, nil
	}

	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	numField := elemType.NumField()
	for i := 0; i < numField; i++ {
		aField := elemType.Field(i)
		if FieldChecker(aField, field) {
			return xunsafe.FieldByIndex(elemType, i), nil
		}
	}

	return nil, fmt.Errorf("not found field %v at struct %v", field, elemType.String())
}
