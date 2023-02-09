package est

import (
	"github.com/viant/xunsafe"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultPkg         = "github.com/viant/velty/est"
	namespace          = "_velty_"
	anonymousNamespace = "Velty_"
)

//Type represents scope type
type Type struct {
	reflect.Type
	fieldsIndex     map[string]int
	fields          []reflect.StructField
	xFields         []*xunsafe.Field
	transients      *int
	anonymousFields map[string]string
}

func (t *Type) AddField(id string, name string, rType reflect.Type) reflect.StructField {
	return t.addField(id, name, rType, false, "")
}

func (t *Type) AddFieldWithTag(id string, name, tag string, rType reflect.Type) reflect.StructField {
	return t.addField(id, name, rType, false, tag)
}

func (t *Type) ValueAccessors() []*xunsafe.Field {
	return t.xFields
}

func (t *Type) EmbedType(rType reflect.Type) reflect.StructField {
	id := t.typeName(rType)

	field := t.addField(id, id, rType, true, "")
	return field
}

func (t *Type) typeName(rType reflect.Type) string {
	var id string

	for rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}

	if rType.Name() != "" {
		idSegments := strings.Split(rType.String(), ".")
		id = idSegments[len(idSegments)-1]
	}

	return id
}

func (t *Type) addField(id string, name string, rType reflect.Type, anonymous bool, defaultTag string) reflect.StructField {
	pkg := ""
	if name == "" {
		if !anonymous {
			name = t.ReserveNewName()
		} else {
			name = t.reserveNewName(anonymousNamespace)
			t.anonymousFields[rType.String()] = name
			id = name
		}
	}

	if name[0] > 'Z' && !anonymous {
		pkg = defaultPkg
	}

	field := reflect.StructField{Name: name, Type: rType, PkgPath: pkg, Anonymous: anonymous}
	t.fields = append(t.fields, field)
	t.Type = reflect.StructOf(t.fields)

	field = t.Type.Field(len(t.fields) - 1)
	t.fields[len(t.fields)-1] = field
	t.xFields = append(t.xFields, xunsafe.NewField(field))
	t.fieldsIndex[id] = len(t.fields) - 1
	return field
}

func (t *Type) ValueAccessor(id string) (*xunsafe.Field, bool) {
	index, found := t.fieldsIndex[id]
	if !found {
		return nil, false
	}

	return t.xFields[index], true
}

func (t *Type) Snapshot() *Type {
	xFields := make([]*xunsafe.Field, len(t.xFields))
	fields := make([]reflect.StructField, len(t.fields))
	types := map[string]int{}
	anonymousFields := map[string]string{}
	copy(xFields, t.xFields)
	copy(fields, t.fields)
	copy(fields, t.fields)

	for i, field := range fields {
		types[field.Name] = i
	}

	for key := range t.anonymousFields {
		anonymousFields[key] = t.anonymousFields[key]
	}

	return &Type{
		Type:            t.Type,
		fieldsIndex:     types,
		fields:          fields,
		xFields:         xFields,
		transients:      t.transients,
		anonymousFields: anonymousFields,
	}
}

func (t *Type) ReserveNewName() string {
	return t.reserveNewName(namespace)
}

func (t *Type) reserveNewName(ns string) string {
	name := ns + "T" + strconv.Itoa(*t.transients)
	*t.transients++
	return name
}

func (t *Type) AnonymousHolder(rType reflect.Type) (string, bool) {
	stringified := rType.String()
	holder, ok := t.anonymousFields[stringified]
	return holder, ok
}

func NewType() *Type {
	transients := 0
	return &Type{
		Type:            reflect.StructOf([]reflect.StructField{}),
		fieldsIndex:     map[string]int{},
		transients:      &transients,
		anonymousFields: map[string]string{},
	}
}
