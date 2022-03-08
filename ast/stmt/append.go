package stmt

import (
	"fmt"
	"github.com/viant/velty/ast"
)

type Append struct {
	Append string
}

func (a Append) AddStatement(statement ast.Statement) error {
	return fmt.Errorf("unexpected statement at Append")
}

func NewAppend(value string) *Append {
	return &Append{Append: value}
}
