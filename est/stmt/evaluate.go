package stmt

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/cache"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/parser"
	"unsafe"
)

type evaluator struct {
	x        *op.Operand
	cache    *cache.Cache
	compiler est.Compiler
	control  est.Control
}

func (e *evaluator) compute(state *est.State) unsafe.Pointer {
	varName := *(*string)(e.x.Exec(state))
	if expr, ok := e.cache.Expression(varName); ok {
		return expr(state)
	}

	block, err := parser.Parse([]byte(varName))
	//TODO: Handler error
	if err != nil {
		panic(err)
	}
	newCompute, err := e.compiler.CompileStmt(block)
	if err != nil {
		panic(err)
	}

	compute, err := newCompute(e.control)
	if err != nil {
		panic(err)
	}

	e.cache.Put(varName, compute)
	return compute(state)
}

func Evaluate(expr *op.Expression, cache *cache.Cache, compiler est.Compiler) (est.New, error) {
	return func(control est.Control) (est.Compute, error) {
		x, err := expr.Operand(control)
		if err != nil {
			return nil, err
		}

		return (&evaluator{
			x:        x,
			cache:    cache,
			compiler: compiler,
			control:  control,
		}).compute, nil
	}, nil
}
