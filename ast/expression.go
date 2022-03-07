package ast

import "reflect"

type Expression interface {
	Type() reflect.Type
}
