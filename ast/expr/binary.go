package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

//Binary represents binary expressions
type Binary struct {
	X     ast.Expression //left operand
	Token ast.Token
	Y     ast.Expression //left operand
}

func (b *Binary) Type() reflect.Type {
	switch b.Token {
	case ast.LEQ, ast.LSS, ast.GTR, ast.GTE, ast.NEQ, ast.EQ:
		return reflect.TypeOf(true)
	}

	if xType := b.X.Type(); xType != nil {
		return xType
	}
	return b.Y.Type()
}

//BinaryExpression creates new *Binary
func BinaryExpression(left ast.Expression, token ast.Token, right ast.Expression) *Binary {
	return &Binary{
		X:     left,
		Token: token,
		Y:     right,
	}
}
