package stmt

import (
	"github.com/viant/velty/internal/ast"
	"github.com/viant/velty/internal/ast/expr"
)

//ForLoop represents regular for loop
type ForLoop struct {
	Init ast.Statement
	Cond ast.Expression
	Body Block
	Post ast.Statement
}

func (r *ForLoop) Statements() []ast.Statement {
	return r.Body.Statements()
}

func (r *ForLoop) AddStatement(statement ast.Statement) {
	r.Body.AddStatement(statement)
}

//ForEach represents for each loop
type ForEach struct {
	Index *expr.Select
	Item  *expr.Select
	Set   ast.Expression
	Body  Block
}

func (f *ForEach) Statements() []ast.Statement {
	return f.Body.Statements()
}

func (f *ForEach) AddStatement(statement ast.Statement) {
	f.Body.AddStatement(statement)
}
