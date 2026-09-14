package velty

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/viant/velty/est/op"
	"github.com/viant/xunsafe"
)

// selectField binds one authored segment. Registration may stop at a recursive
// type, while each expression retains its own finite field/index/call chain.
func (p *Planner) selectField(parent *op.Selector, rType reflect.Type, name string) (*op.Selector, error) {
	path := p.selectorFieldPath(rType, name, nil)
	if len(path) == 0 {
		return nil, fmt.Errorf("not found field %v at %v", name, rType)
	}
	id := parent.ID + fieldSeparator + name
	for _, field := range path {
		offset := uintptr(0)
		if parent.Field != nil {
			offset = parent.ParentOffset + parent.Field.Offset
		}
		indirect := parent.Indirect || field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Slice
		parent = op.SelectorWithField(id, xunsafe.NewField(field), parent, indirect, offset)
		parent.NodeID = p.nextNodeID()
	}
	return parent, nil
}

// selectorFieldPath applies the same aliases, omissions and anonymous prefixes
// as registration, retaining embedded pointer steps instead of flattening offsets.
func (p *Planner) selectorFieldPath(rType reflect.Type, name string, ancestors *CycleDetector) []reflect.StructField {
	rType, _ = elemIfNeeded(rType)
	if rType.Kind() != reflect.Struct || rType == TimeType || name == "_" {
		return nil
	}
	next, cycle := ancestors.Child(rType, nil)
	if cycle {
		return nil
	}
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		tag := Parse(field.Tag.Get(velty))
		if tag.Omit || field.Anonymous {
			continue
		}
		if (len(tag.Names) == 0 && field.Name == name) || (name != "_" && tag.nameEqual(name)) {
			return []reflect.StructField{field}
		}
	}
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		tag := Parse(field.Tag.Get(velty))
		if !field.Anonymous || tag.Omit || !strings.HasPrefix(name, tag.Prefix) {
			continue
		}
		if path := p.selectorFieldPath(field.Type, strings.TrimPrefix(name, tag.Prefix), next); len(path) != 0 {
			return append([]reflect.StructField{field}, path...)
		}
	}
	return nil
}
