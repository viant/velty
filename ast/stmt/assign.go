package stmt

import (
	ast2 "github.com/viant/velty/ast"
)

//Statement represents assign statement i.e. $var = 10
type Statement struct {
	X  ast2.Expression //left operand
	Op ast2.Operand    // =
	Y  ast2.Expression //right operand
}
