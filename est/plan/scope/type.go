package scope

import (
	"github.com/viant/xunsafe"
	"reflect"
)

const defaultPkg = "github.com/viant/vetly"

//Type represents scope type
type Type struct {
	reflect.Type
	fields []reflect.StructField
}

func (t *Type) AddField(name string, fType reflect.Type) *xunsafe.Field {
	pkg := ""
	if name[0] > 'Z' {
		pkg = defaultPkg
	}
	if fType.Kind() == reflect.Ptr {
		fType = fType.Elem()
	}
	idx := len(t.fields)
	t.fields = append(t.fields, reflect.StructField{Name: name, Type: fType, PkgPath: pkg})
	t.Type = reflect.StructOf(t.fields)
	result := xunsafe.NewField(t.Type.Field(idx))
	return result
}

func NewType() *Type {
	return &Type{Type: reflect.StructOf([]reflect.StructField{})}
}
