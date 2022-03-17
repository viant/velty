package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

type Select struct {
	ID string
	X  ast.Expression
}

func (s Select) Type() reflect.Type {
	return nil
}
