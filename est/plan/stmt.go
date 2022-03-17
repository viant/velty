package plan

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/est"
	estmt "github.com/viant/velty/est/stmt"
	"unsafe"
)

func (p *Planner) compileStmt(statement ast.Statement) (est.New, error) {

	switch actual := statement.(type) {
	case *stmt.Statement:
		return p.computeDirectAssignment(actual)
	case *stmt.Append:
		return p.compileAppend(actual)
	case *expr.Select:
		return p.compileStmtSelector(actual)
	case *stmt.Block:
		return p.compileStmt(actual.Stmt)
	case *stmt.If:
		return p.compileIf(actual)
	case *stmt.Range:
		return p.compileForLoop(actual)
	case *stmt.ForEach:
		return p.compileForEachLoop(actual)
	case []ast.Statement:
		return p.compileBlock(&stmt.Block{Stmt: actual})
	}

	return nil, fmt.Errorf("unsupported stmt: %T", statement)
}

func (p *Planner) computeDirectAssignment(actual *stmt.Statement) (est.New, error) {
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
}

func (p *Planner) compileIf(actual *stmt.If) (est.New, error) {
	cond, err := p.compileExpr(actual.Condition)
	if err != nil {
		return nil, err
	}

	body, err := p.compileStmt(&actual.Body)
	if err != nil {
		return nil, err
	}

	var elseIf est.New
	if actual.Else != nil {
		elseIf, err = p.compileStmt(actual.Else)
		if err != nil {
			return nil, err
		}
	}

	return estmt.NewIf(cond, body, elseIf)
}

func (p *Planner) compileForLoop(actual *stmt.Range) (est.New, error) {
	init, err := p.compileStmt(actual.Init)

	if err != nil {
		return nil, err
	}

	post, err := p.compileStmt(actual.Post)
	if err != nil {
		return nil, err
	}

	condition, err := p.compileExpr(actual.Cond)
	if err != nil {
		return nil, err
	}

	block, err := p.newCompute(&actual.Body)
	if err != nil {
		return nil, err
	}

	return estmt.ForLoop(init, post, condition, block)
}

func (p *Planner) compileForEachLoop(actual *stmt.ForEach) (est.New, error) {

	sliceSelector, err := p.compileExpr(actual.Set)
	if err != nil {
		return nil, err
	}

	itemType := sliceSelector.Type.Elem()

	if err := p.DefineVariable(actual.Item.ID, itemType); err != nil {
		return nil, err
	}

	selector, err := p.compileExpr(actual.Item)
	if err != nil {
		return nil, err
	}

	block, err := p.compileBlock(&actual.Body)
	if err != nil {
		return nil, err
	}
	return estmt.ForEachLoop(block, selector, sliceSelector)
}

func (p *Planner) compileAppend(actual *stmt.Append) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		ptr := unsafe.Pointer(&actual.Append)
		return func(mem *est.State) unsafe.Pointer {
			mem.Buffer.AppendString(actual.Append)
			return ptr
		}, nil
	}, nil
}
