package functions

import (
	"reflect"
	"strings"
)

type (
	Registry struct {
		Entries map[string]*Entry
	}

	Entry struct {
		Metadata interface{}
		Handler  interface{}
	}

	FunctionKindMetadata struct {
		kinds []reflect.Kind
	}

	FunctionNamespaceMetadata struct {
		rType reflect.Type
	}
)

func CopyInstance() *Registry {
	registry := NewRegistry()
	for key, entry := range registryInstance.Entries {
		anEntry := &Entry{
			Metadata: entry.Metadata,
			Handler:  entry.Handler,
		}

		registry.Entries[key] = anEntry
	}

	return registry
}

func NewRegistry() *Registry {
	return &Registry{
		Entries: map[string]*Entry{},
	}
}

var registryInstance = NewRegistry()

func NewFunctionNamespace(rType reflect.Type) *FunctionNamespaceMetadata {
	return &FunctionNamespaceMetadata{rType: rType}
}

func NewFunctionKind(kinds []reflect.Kind) *FunctionKindMetadata {
	return &FunctionKindMetadata{
		kinds: kinds,
	}
}

func NewEntry(handler interface{}, metadata interface{}) *Entry {
	return &Entry{
		Metadata: handler,
		Handler:  metadata,
	}
}

func (f *FunctionKindMetadata) Kinds() []reflect.Kind {
	result := make([]reflect.Kind, len(f.kinds))
	copy(result, f.kinds)
	return result
}

func (f *FunctionNamespaceMetadata) Type() reflect.Type {
	return f.rType
}

func (r *Registry) DefineNs(ns string, def *Entry) string {
	r.Entries[strings.Trim(ns, "${}")] = def
	return ns
}

func (r *Registry) IsDefined(ns string) bool {
	_, ok := r.Entries[ns]
	return ok
}

var interfaceType reflect.Type

func init() {
	type foo struct {
		aField interface{}
	}

	interfaceType = reflect.ValueOf(foo{}).Field(0).Type()
}
