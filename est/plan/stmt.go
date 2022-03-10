package plan

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/est"
	estmt "github.com/viant/velty/est/stmt"
)

func (p *Planner) compileStmt(statement ast.Statement) (est.New, error) {

	switch actual := statement.(type) {
	case *stmt.Statement:
		x, err := p.compileExpr(actual.X)
		if err != nil {
			return nil, err
		}
		y, err := p.compileExpr(actual.Y)
		if err != nil {
			return nil, err
		}
		if err = p.adjustSelector(x, y.Type); err != nil {
			return nil, err
		}
		return estmt.Assign(x, y)
	case *stmt.Append:
		fmt.Printf("append: !%v!\n", actual.Append)
	case *expr.Select:
		return p.compileStmtSelector(actual)
	case *stmt.Block:
		return p.compileStmt(actual.Stmt)
	}
	return nil, fmt.Errorf("unsupported stmt: %T", statement)
}
