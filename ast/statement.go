package ast

type Statement interface {
}

type StatementContainer interface {
	AddStatement(statement Statement)
	Statements() []Statement
}
