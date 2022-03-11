package expr

import (
	"github.com/viant/velty/ast"
	"reflect"
)

type Parentheses struct {
	Parentheses ast.Expression
}

func (p *Parentheses) Type() reflect.Type {
	return p.Parentheses.Type()
}
