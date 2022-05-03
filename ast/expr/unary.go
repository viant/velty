package expr

import (
	ast2 "github.com/viant/velty/ast"
	"reflect"
)

//Unary represents unary expression
type Unary struct {
	Token ast2.Token
	X     ast2.Expression
}

func (u *Unary) Type() reflect.Type {
	return u.X.Type()
}
