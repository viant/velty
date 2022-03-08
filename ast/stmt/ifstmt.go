package stmt

import "github.com/viant/velty/ast"

type If struct {
	Condition ast.Expression
	Body      Block
	Else      *If
}

func (i *If) AddStatement(statement ast.Statement) error {
	return i.Body.AddStatement(statement)
}

func (i *If) AddCondition(next *If) {
	temp := i
	for temp.Else != nil {
		temp = temp.Else
	}

	temp.Else = next
}

type Condition interface {
	AddCondition(condition *If)
}
