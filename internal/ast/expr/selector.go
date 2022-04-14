package expr

import (
	"github.com/viant/velty/internal/ast"
	"reflect"
)

//Select represents dynamic variable
type Select struct {
	ID       string
	X        ast.Expression
	FullName string
}

func (s Select) Type() reflect.Type {
	return nil
}
