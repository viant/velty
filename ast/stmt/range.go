package stmt

import (
	ast2 "github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
)

//ForLoop represents regular for loop
type ForLoop struct {
	Init ast2.Statement
	Cond ast2.Expression
	Body Block
	Post ast2.Statement
}

func (r *ForLoop) Statements() []ast2.Statement {
	return r.Body.Statements()
}

func (r *ForLoop) AddStatement(statement ast2.Statement) {
	r.Body.AddStatement(statement)
}

//ForEach represents for each loop
type ForEach struct {
	Index *expr.Select
	Item  *expr.Select
	Set   ast2.Expression
	Body  Block
}

func (f *ForEach) Statements() []ast2.Statement {
	return f.Body.Statements()
}

func (f *ForEach) AddStatement(statement ast2.Statement) {
	f.Body.AddStatement(statement)
}
