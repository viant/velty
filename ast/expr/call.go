package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

type Call struct {
	Args []ast.Expression
	X    ast.Expression
}

func (c *Call) Type() reflect.Type {
	return nil
}

type SliceIndex struct {
	X ast.Expression
	Y ast.Expression
}

func (s *SliceIndex) Type() reflect.Type {
	return reflect.TypeOf(0)
}
