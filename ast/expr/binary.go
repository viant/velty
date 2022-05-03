package expr

import (
	ast2 "github.com/viant/velty/ast"
	"reflect"
)

//Binary represents binary expressions
type Binary struct {
	X     ast2.Expression //left operand
	Token ast2.Token
	Y     ast2.Expression //left operand
}

func (b *Binary) Type() reflect.Type {
	switch b.Token {
	case ast2.LEQ, ast2.LSS, ast2.GTR, ast2.GTE, ast2.NEQ, ast2.EQ:
		return reflect.TypeOf(true)
	}

	if xType := b.X.Type(); xType != nil {
		return xType
	}
	return b.Y.Type()
}

//BinaryExpression creates new *Binary
func BinaryExpression(left ast2.Expression, token ast2.Token, right ast2.Expression) *Binary {
	return &Binary{
		X:     left,
		Token: token,
		Y:     right,
	}
}
