package functions

import "reflect"

//TypeParser parses type string representation into reflect.Type
type TypeParser func(typeRepresentation string) (reflect.Type, error)
