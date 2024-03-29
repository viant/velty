package functions

import "reflect"

var (
	stringType = reflect.TypeOf("")
	boolType   = reflect.TypeOf(false)
	floatType  = reflect.TypeOf(0.0)
)

type Types struct{}

func (t Types) IsInt(value interface{}) bool {
	_, ok := value.(int)
	return ok
}

func (t Types) IsFloat64(value interface{}) bool {
	_, ok := value.(float64)
	return ok
}

func (t Types) IsString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func (t Types) IsBool(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}
