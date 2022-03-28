package expr

import (
	"github.com/viant/velty/internal/ast"
	"reflect"
)

//Call represents function call
type Call struct {
	Args []ast.Expression
	X    ast.Expression
}

func (c *Call) Type() reflect.Type {
	return nil
}

//SliceIndex represents slice accessor
type SliceIndex struct {
	X ast.Expression
	Y ast.Expression
}

func (s *SliceIndex) Type() reflect.Type {
	return reflect.TypeOf(0)
}
