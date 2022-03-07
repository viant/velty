package stmt

import "github.com/viant/velty/ast"

type If struct {
	Condition ast.Expression
	Body      Block
	Else      *If
}
