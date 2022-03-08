package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

type Binary struct {
	X     ast.Expression //left operand
	Token ast.Token
	Y     ast.Expression //left operand
}

func (b *Binary) Type() reflect.Type {
	if xType := b.X.Type(); xType != nil {
		return xType
	}
	return b.Y.Type()
}
