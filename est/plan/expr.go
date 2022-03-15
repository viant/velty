package plan

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est/op"
)

func (p *Planner) compileExpr(e ast.Expression) (*op.Expression, error) {
	switch actual := e.(type) {
	case *expr.Literal:
		return p.literalExpr(actual)
	case *expr.Select:
		return p.selectorExpr(actual)
	case *expr.Binary:
		return p.compileBinary(actual)
	case *expr.Parentheses:
		return p.compileExpr(actual.P)
	}

	return nil, fmt.Errorf("unsupported expr: %T", e)
}
