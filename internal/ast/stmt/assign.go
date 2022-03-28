package stmt

import (
	"github.com/viant/velty/internal/ast"
)

//Statement represents assign statement i.e. $var = 10
type Statement struct {
	X  ast.Expression //left operand
	Op ast.Operand    // =
	Y  ast.Expression //right operand
}
