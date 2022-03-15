package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

type Parentheses struct {
	P ast.Expression
}

func (p *Parentheses) Type() reflect.Type {
	return p.P.Type()
}
