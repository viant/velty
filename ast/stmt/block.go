package stmt

import "github.com/viant/velty/ast"

//Block represents block of statements
type Block struct {
	Stmt []ast.Statement
}

func (b *Block) Statements() []ast.Statement {
	return b.Stmt
}

func (b *Block) AddStatement(statement ast.Statement) {
	b.Stmt = append(b.Stmt, statement)
}
