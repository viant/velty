package functions

import (
	"fmt"
	"reflect"
)

type Maps struct{}

func (m Maps) Has(aMap interface{}, key interface{}) (bool, error) {
	index, err := m.value(aMap, key)
	if err != nil {
		return false, err
	}

	return index.IsValid(), nil
}

func (m Maps) value(aMap interface{}, key interface{}) (reflect.Value, error) {
	if aMap == nil || key == nil {
		return reflect.Value{}, nil
	}

	valueOf := reflect.ValueOf(aMap)
	if valueOf.Kind() != reflect.Map {
		return reflect.Value{}, fmt.Errorf("incorrect arg type, wanted Map, got %T", aMap)
	}

	index := valueOf.MapIndex(reflect.ValueOf(key))
	return index, nil
}

func (m Maps) Get(aMap interface{}, key interface{}) (interface{}, error) {
	value, err := m.value(aMap, key)
	if err != nil {
		return nil, err
	}

	if !value.IsValid() {
		return nil, nil
	}

	return value.Interface(), nil
}

func (m Maps) GetInt(aMap interface{}, key interface{}) (int, error) {
	iface, err := m.Get(aMap, key)
	if err != nil {
		return 0, err
	}

	i, ok := iface.(int)
	if !ok {
		return 0, fmt.Errorf("unexpected map value type, watned %T got %T", i, iface)
	}

	return i, nil
}

func (m Maps) GetFloat(aMap interface{}, key interface{}) (float64, error) {
	iface, err := m.Get(aMap, key)
	if err != nil {
		return 0, err
	}

	i, ok := iface.(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected map value type, watned %T got %T", i, iface)
	}

	return i, nil
}

func (m Maps) GetBool(aMap interface{}, key interface{}) (bool, error) {
	iface, err := m.Get(aMap, key)
	if err != nil {
		return false, err
	}

	i, ok := iface.(bool)
	if !ok {
		return false, fmt.Errorf("unexpected map value type, watned %T got %T", i, iface)
	}

	return i, nil
}

func (m Maps) GetString(aMap interface{}, key interface{}) (string, error) {
	iface, err := m.Get(aMap, key)
	if err != nil {
		return "", err
	}

	i, ok := iface.(string)
	if !ok {
		return "", fmt.Errorf("unexpected map value type, watned %T got %T", i, iface)
	}

	return i, nil
}

var HasKeyFunc = &StaticKindFunc{
	kind: reflect.Map,
	handler: func(args ...interface{}) (bool, error) {
		if len(args) != 2 {
			return false, fmt.Errorf("unexpected number of args, expected %v got %v", 2, len(args))
		}

		return reflect.ValueOf(args[0]).MapIndex(reflect.ValueOf(args[1])).IsValid(), nil
	},
	resultType: reflect.TypeOf(true),
}
