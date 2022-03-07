package stmt

import "github.com/viant/velty/ast"

type Range struct {
	Init ast.Statement
	Cond ast.Expression
	Body Block
	Post ast.Statement
}
