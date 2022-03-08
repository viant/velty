package ast

type Statement interface {
	AddStatement(statement Statement) error
}
