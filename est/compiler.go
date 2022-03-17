package est

import "github.com/viant/velty/ast"

type Compiler interface {
	CompileStmt(statement ast.Statement) (New, error)
}
