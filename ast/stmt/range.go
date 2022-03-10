package stmt

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
)

type Range struct {
	Init ast.Statement
	Cond ast.Expression
	Body Block
	Post ast.Statement
}

func (r *Range) AddStatement(statement ast.Statement) {
	r.Body.AddStatement(statement)
}

type ForEach struct {
	Index *expr.Select
	Item  *expr.Select
	Set   *expr.Select
	Body  Block
}

func (f *ForEach) AddStatement(statement ast.Statement) {
	f.Body.AddStatement(statement)
}
