package stmt

import "github.com/viant/velty/ast"

type If struct {
	Condition ast.Expression
	Body      Block
	Else      *If
}

func (i *If) Statements() []ast.Statement {
	return i.Body.Statements()
}

func (i *If) AddStatement(statement ast.Statement) {
	temp := i
	for temp.Else != nil {
		temp = temp.Else
	}
	temp.Body.AddStatement(statement)
}

func (i *If) AddCondition(next *If) {
	temp := i
	for temp.Else != nil {
		temp = temp.Else
	}

	temp.Else = next
}

type ConditionContainer interface {
	AddCondition(condition *If)
}
