package velty

import (
	"github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/est"
	estmt "github.com/viant/velty/est/stmt"
)

func (p *Planner) compileBlock(root *stmt.Block) (est.New, error) {
	var newComputers = make([]est.New, len(root.Stmt))
	var err error
	for i, item := range root.Stmt {
		if newComputers[i], err = p.compileStmt(item); err != nil {
			return nil, err
		}
	}
	return estmt.NewBlock(newComputers), nil
}
