package est

import (
	"github.com/viant/xunsafe"
	"reflect"
	"strings"
)

const defaultPkg = "github.com/viant/velty/est"

//Type represents scope type
type Type struct {
	reflect.Type
	types   map[string]int
	fields  []reflect.StructField
	xFields []*xunsafe.Field
}

func (t *Type) AddField(id string, name string, rType reflect.Type) reflect.StructField {
	return t.addField(id, name, rType, false)
}

func (t *Type) ValueAccessors() []*xunsafe.Field {
	return t.xFields
}

func (t *Type) EmbedType(rType reflect.Type) reflect.StructField {
	idSegments := strings.Split(rType.String(), ".")
	id := idSegments[len(idSegments)-1]

	field := t.addField(id, id, rType, true)
	return field
}

func (t *Type) addField(id string, name string, rType reflect.Type, anonymous bool) reflect.StructField {
	pkg := ""
	if name[0] > 'Z' {
		pkg = defaultPkg
	}

	field := reflect.StructField{Name: name, Type: rType, PkgPath: pkg, Anonymous: anonymous}
	t.fields = append(t.fields, field)
	t.Type = reflect.StructOf(t.fields)

	field = t.Type.Field(len(t.fields) - 1)
	t.fields[len(t.fields)-1] = field
	t.xFields = append(t.xFields, xunsafe.NewField(field))
	t.types[id] = len(t.fields) - 1
	return field
}

func (t *Type) ValueAccessor(id string) (*xunsafe.Field, bool) {
	index, found := t.types[id]
	if !found {
		return nil, false
	}

	return t.xFields[index], true
}

func (t *Type) Snapshot() *Type {
	xFields := make([]*xunsafe.Field, len(t.xFields))
	fields := make([]reflect.StructField, len(t.fields))
	types := map[string]int{}
	copy(xFields, t.xFields)
	copy(fields, t.fields)
	copy(fields, t.fields)

	for i, field := range fields {
		types[field.Name] = i
	}

	return &Type{
		Type:    t.Type,
		types:   types,
		fields:  fields,
		xFields: xFields,
	}
}

func NewType() *Type {
	return &Type{
		Type:  reflect.StructOf([]reflect.StructField{}),
		types: map[string]int{},
	}
}
