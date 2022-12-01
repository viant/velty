package velty

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	stmt2 "github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/stmt"
	"github.com/viant/velty/est/stmt/assign"
	"github.com/viant/xunsafe/converter"
	"unsafe"
)

func (p *Planner) compileStmt(statement ast.Statement) (est.New, error) {
	switch actual := statement.(type) {
	case *stmt2.Statement:
		return p.computeAssignment(actual)
	case *stmt2.Append:
		return p.compileAppend(actual)
	case *expr.Select:
		return p.compileStmtSelector(actual)
	case *stmt2.Block:
		return p.compileStmt(actual.Stmt)
	case *stmt2.If:
		return p.compileIf(actual)
	case *stmt2.ForLoop:
		return p.compileForLoop(actual)
	case *stmt2.ForEach:
		return p.compileForEachLoop(actual)
	case []ast.Statement:
		return p.compileBlock(&stmt2.Block{Stmt: actual})
	case *stmt2.Evaluate:
		return p.compileEvaluate(actual)
	}

	return nil, fmt.Errorf("unsupported stmt: %T", statement)
}

func (p *Planner) computeAssignment(actual *stmt2.Statement) (est.New, error) {
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

	unify, err := converter.Unify(x.Type, y.Type)
	if err != nil {
		return nil, err
	}

	y.Unify = unify.Y
	y.Type = unify.RType

	return assign.Assign(x, y)
}

func (p *Planner) compileIf(actual *stmt2.If) (est.New, error) {
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

	return stmt.NewIf(cond, body, elseIf)
}

func (p *Planner) compileForLoop(actual *stmt2.ForLoop) (est.New, error) {
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

	return stmt.ForLoop(init, post, condition, block)
}

func (p *Planner) compileForEachLoop(actual *stmt2.ForEach) (est.New, error) {
	sliceSelector, err := p.compileExpr(actual.Set)
	if err != nil {
		return nil, err
	}

	if sliceSelector.Type == nil {
		return nop(), nil
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
	return stmt.ForEachLoop(block, selector, sliceSelector)
}

func nop() est.New {
	return func(control est.Control) (est.Compute, error) {
		return func(state *est.State) unsafe.Pointer {
			return est.EmptyStringPtr
		}, nil
	}
}

func (p *Planner) compileAppend(actual *stmt2.Append) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		ptr := unsafe.Pointer(&actual.Append)
		return func(state *est.State) unsafe.Pointer {
			state.Buffer.AppendStringWithoutEscaping(actual.Append)
			return ptr
		}, nil
	}, nil
}

func (p *Planner) compileEvaluate(actual *stmt2.Evaluate) (est.New, error) {
	selector, err := p.compileExpr(actual.X)
	if err != nil {
		return nil, err
	}

	return evaluate(selector, p.cache, p)
}
