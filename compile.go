package velty

import (
	"context"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/parser"
	"github.com/viant/xunsafe"
	"reflect"
)

// Compile create Execution Plan and State provider for the Execution Plan.
func (p *Planner) Compile(template []byte) (*est.Execution, func() *est.State, error) {
	var (
		root *stmt.Block
		err  error
	)

	hasParserHooks := p.listener != nil || p.adjuster != nil

	if !hasParserHooks {
		// Fast path: parse without spans and compile directly.
		root, err = parser.Parse(template)
		if err != nil {
			return nil, nil, err
		}
	} else {
		// Instrumented path: parse with spans and run listener/adjuster hooks.
		var spans map[ast.Node]parser.NodeSpan
		root, spans, err = parser.ParseWithSpansDetailed(template)
		if err != nil {
			return nil, nil, err
		}

		// Build symbol seeds for resolution helpers
		params := make([]string, 0)
		for _, sel := range p.selectors.Selectors() {
			if sel.Parent == nil && sel.IsFieldSelector {
				params = append(params, sel.ID)
			}
		}
		namespaces := []string{}
		funcs := []string{}
		if p.Functions != nil {
			namespaces = p.Functions.NamespaceNames()
			funcs = p.Functions.StandaloneNames()
		}
		seeds := &SymbolSeeds{Params: params, Namespaces: namespaces, Standalone: funcs}
		// apply hooks with evaluate config and seeds; ignore patches here (integrator may use them separately)
		hookCtx := &ParserContext{Scope: NewScope(nil)}
		hookCtx.InitSource("", template)
		hookCtx.EvalConfig = p.evalConfig
		if transformed, _, err := applyParserHooksWithConfigAndSpans(hookCtx, root, p.listener, p.adjuster, spans, seeds); err == nil && transformed != nil {
			root = transformed
		} else if err != nil {
			return nil, nil, err
		}
	}

	exec, err := p.newExecution(root)
	if err != nil {
		return nil, nil, err
	}

	newState := p.stateProvider()
	return exec, newState, err
}

func (p *Planner) stateProvider() func() *est.State {
	return func() *est.State {
		mem := reflect.New(p.Type.Type).Interface()
		state := &est.State{
			Mem:          mem,
			MemPtr:       xunsafe.AsPointer(mem),
			Buffer:       est.NewBuffer(p.bufferSize, p.escapeHTML),
			StateType:    p.Type,
			PanicOnError: p.panicOnError,
			Ctx:          context.Background(),
		}

		return state
	}
}

func (p *Planner) newExecution(root *stmt.Block) (*est.Execution, error) {
	compute, err := p.newCompute(root)
	if err != nil {
		return nil, err
	}

	exec := est.NewExecution(compute)
	exec.PanicOnError = p.panicOnError
	return exec, nil
}

func (p *Planner) newCompute(root *stmt.Block) (est.Compute, error) {
	computeNew, err := p.compileBlock(root)
	if err != nil {
		return nil, err
	}
	compute, err := computeNew(*p.Control)
	if err != nil {
		return nil, err
	}
	return compute, nil
}
