package stmt

import (
	"github.com/viant/igo/exec/est"
	"github.com/viant/velty/ast"
)

type Statement struct {
	X  ast.Expression //left operand
	Op est.Operand    // =
	Y  ast.Expression //right operand
}
