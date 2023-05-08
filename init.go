package velty

import (
	"github.com/viant/velty/functions"
	"reflect"
)

func init() {
	functions.FieldChecker = func(field reflect.StructField, fieldName string) bool {
		aTag := Parse(field.Tag.Get("velty"))
		for _, name := range aTag.Names {
			if name == fieldName {
				return true
			}
		}

		return false
	}
}
