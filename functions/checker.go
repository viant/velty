package functions

import "reflect"

type FieldNameChecker func(field reflect.StructField, fieldName string) bool

var FieldChecker FieldNameChecker = func(field reflect.StructField, fieldName string) bool {
	return false
}
