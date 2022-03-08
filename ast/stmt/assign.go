package stmt

import (
	"fmt"
	"github.com/viant/velty/ast"
)

type Statement struct {
	X  ast.Expression //left operand
	Op ast.Operand    // =
	Y  ast.Expression //right operand
}

func (s *Statement) AddStatement(_ ast.Statement) error {
	return fmt.Errorf("unexpected add statement for type Statement")
}
