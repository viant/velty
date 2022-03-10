package stmt

import (
	"github.com/viant/velty/ast"
)

type Statement struct {
	X  ast.Expression //left operand
	Op ast.Operand    // =
	Y  Assignable     //right operand
}

type Assignable interface {
}
