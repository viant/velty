package stmt

import "github.com/viant/velty/ast"

type Range struct {
	Init ast.Statement
	Cond ast.Expression
	Body Block
	Post ast.Statement
}

func (r *Range) AddStatement(statement ast.Statement) error {
	return r.Body.AddStatement(statement)
}
