package expr

import "github.com/viant/velty/ast"

type Unary struct {
	Token ast.Token
	X     ast.Expression
}
