package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

//Parentheses used to add precedence to the expression over other expression
type Parentheses struct {
	P ast.Expression
}

func (p *Parentheses) Type() reflect.Type {
	return p.P.Type()
}
