package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

type Binary struct {
	X     ast.Expression //left operand
	Token ast.Token
	Y     ast.Expression //left operand
	t     reflect.Type
}

func (b *Binary) Type() reflect.Type {
	return b.t
}
