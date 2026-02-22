package velty

import (
	"github.com/viant/velty/ast"
	aexpr "github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
	"github.com/viant/velty/parser"
	"reflect"
)

// applyParserHooks traverses the AST, firing listener events and applying adjuster transformations.
func applyParserHooks(root *stmt.Block, listener ParserListener, adjuster NodeAdjuster) (*stmt.Block, error) {
	ctx := &ParserContext{Scope: NewScope(nil)}
	for i, s := range root.Stmt {
		node, err := walkNode(s, listener, adjuster, ctx)
		if err != nil {
			return nil, err
		}
		if st, ok := node.(ast.Statement); ok {
			root.Stmt[i] = st
		}
	}
	return root, nil
}

// applyParserHooksWithSource initializes context source mapping and then traverses.
func applyParserHooksWithSource(file string, src []byte, root *stmt.Block, listener ParserListener, adjuster NodeAdjuster) (*stmt.Block, error) {
	ctx := &ParserContext{Scope: NewScope(nil)}
	ctx.InitSource(file, src)
	spans := spansFromSource(src)
	return applyParserHooksWithSourceAndSpans(ctx, root, listener, adjuster, spans)
}

func applyParserHooksWithSourceAndSpans(ctx *ParserContext, root *stmt.Block, listener ParserListener, adjuster NodeAdjuster, spans map[ast.Node]parser.NodeSpan) (*stmt.Block, error) {
	// seed spans recorded by the parser into this context
	for n, s := range spans {
		ctx.SetSpan(n, Span{Start: s.Start, End: s.End})
	}
	for i, s := range root.Stmt {
		node, err := walkNode(s, listener, adjuster, ctx)
		if err != nil {
			return nil, err
		}
		if st, ok := node.(ast.Statement); ok {
			root.Stmt[i] = st
		}
	}
	return root, nil
}

// applyParserHooksWithConfig initializes context with source and evaluate config, then traverses.
type SymbolSeeds struct {
	Params     []string
	Namespaces []string
	Standalone []string
	Consts     []string
}

func applyParserHooksWithConfig(file string, src []byte, root *stmt.Block, listener ParserListener, adjuster NodeAdjuster, eval EvaluateConfig, seeds ...*SymbolSeeds) (*stmt.Block, *ParserContext, error) {
	ctx := &ParserContext{Scope: NewScope(nil)}
	ctx.InitSource(file, src)
	ctx.EvalConfig = eval
	spans := spansFromSource(src)
	return applyParserHooksWithConfigAndSpans(ctx, root, listener, adjuster, spans, seeds...)
}

func spansFromSource(src []byte) map[ast.Node]parser.NodeSpan {
	if len(src) == 0 {
		return nil
	}
	_, spans, err := parser.ParseWithSpansDetailed(src)
	if err != nil {
		return nil
	}
	return spans
}

func applyParserHooksWithConfigAndSpans(ctx *ParserContext, root *stmt.Block, listener ParserListener, adjuster NodeAdjuster, spans map[ast.Node]parser.NodeSpan, seeds ...*SymbolSeeds) (*stmt.Block, *ParserContext, error) {
	// seed spans recorded by the parser into this context
	for n, s := range spans {
		ctx.SetSpan(n, Span{Start: s.Start, End: s.End})
	}
	if len(seeds) > 0 && seeds[0] != nil {
		s := seeds[0]
		ctx.SeedSymbols(s.Params, s.Namespaces, s.Standalone, s.Consts)
	}
	for i, s := range root.Stmt {
		node, err := walkNode(s, listener, adjuster, ctx)
		if err != nil {
			return nil, nil, err
		}
		if st, ok := node.(ast.Statement); ok {
			root.Stmt[i] = st
		}
	}
	return root, ctx, nil
}

// walkNode processes a single AST node: fire enter/exit events and apply adjuster.
func walkNode(n ast.Node, listener ParserListener, adjuster NodeAdjuster, ctx *ParserContext) (ast.Node, error) {
	// build event with span/position if available
	enter := Event{Type: EventEnterNode, Node: n, Context: ctx, ExprContext: ctx.CurrentExprContext()}
	if s, ok := ctx.GetSpan(n); ok {
		enter.Span = s
		enter.Position = ctx.ResolvePosition(s)
	}
	if listener != nil {
		listener.OnEvent(enter)
	}
	ctx.PushNode(n)
	if adjuster != nil {
		act, err := adjuster.Adjust(n, ctx)
		if err != nil {
			return nil, err
		}
		switch act.Kind {
		case ActionRemove:
			// do not traverse children; exit event will still fire with nil Node
			ctx.PopNode()
			if listener != nil {
				listener.OnEvent(Event{Type: EventExitNode, Node: n, Context: ctx, ExprContext: ctx.CurrentExprContext()})
			}
			return nil, nil
		case ActionReplace:
			if act.Node != nil {
				n = act.Node
			}
			if act.SkipChildren {
				// skip traversing children
				ctx.PopNode()
				if listener != nil {
					listener.OnEvent(Event{Type: EventExitNode, Node: n, Context: ctx, ExprContext: ctx.CurrentExprContext()})
				}
				return n, nil
			}
		case ActionPatchOnly:
			if act.SkipChildren {
				ctx.PopNode()
				if listener != nil {
					listener.OnEvent(Event{Type: EventExitNode, Node: n, Context: ctx, ExprContext: ctx.CurrentExprContext()})
				}
				return n, nil
			}
		case ActionKeep:
			// continue
		}
	}
	switch node := n.(type) {
	case *stmt.Block:
		for i, child := range node.Stmt {
			newChild, err := walkNode(child, listener, adjuster, ctx)
			if err != nil {
				return nil, err
			}
			if newChild == nil {
				// removal: drop child
				node.Stmt = append(node.Stmt[:i], node.Stmt[i+1:]...)
				i--
				continue
			}
			if st, ok := newChild.(ast.Statement); ok {
				node.Stmt[i] = st
			}
		}
	case *stmt.If:
		// Distinguish top-level if from synthetic else-body wrappers.
		switch ctx.CurrentExprContext().Kind {
		case CtxElseBody:
			// Synthetic If created for #else: we want body selectors to see CtxElseBody.
			for i, child := range node.Body.Stmt {
				newChild, err := walkNode(child, listener, adjuster, ctx)
				if err != nil {
					return nil, err
				}
				if newChild == nil {
					node.Body.Stmt = append(node.Body.Stmt[:i], node.Body.Stmt[i+1:]...)
					i--
					continue
				}
				if st, ok := newChild.(ast.Statement); ok {
					node.Body.Stmt[i] = st
				}
			}
		default:
			// Regular if/elseif chains.
			// If condition
			ctx.PushExprContext(ExprContext{Kind: CtxIfCond, ArgIdx: -1})
			_ = walkExpr(node.Condition, listener, adjuster, ctx)
			ctx.PopExprContext()

			// If body
			ctx.PushExprContext(ExprContext{Kind: CtxIfBody, ArgIdx: -1})
			for i, child := range node.Body.Stmt {
				newChild, err := walkNode(child, listener, adjuster, ctx)
				if err != nil {
					return nil, err
				}
				if newChild == nil {
					node.Body.Stmt = append(node.Body.Stmt[:i], node.Body.Stmt[i+1:]...)
					i--
					continue
				}
				if st, ok := newChild.(ast.Statement); ok {
					node.Body.Stmt[i] = st
				}
			}
			ctx.PopExprContext()

			if node.Else != nil {
				// Else branch can represent either else-if (Condition != nil) or else body
				if node.Else.Condition != nil {
					// Else-if condition
					ctx.PushExprContext(ExprContext{Kind: CtxElseIfCond, ArgIdx: -1})
					_ = walkExpr(node.Else.Condition, listener, adjuster, ctx)
					ctx.PopExprContext()

					// Else-if body
					ctx.PushExprContext(ExprContext{Kind: CtxIfBody, ArgIdx: -1})
					for i, child := range node.Else.Body.Stmt {
						newChild, err := walkNode(child, listener, adjuster, ctx)
						if err != nil {
							return nil, err
						}
						if newChild == nil {
							node.Else.Body.Stmt = append(node.Else.Body.Stmt[:i], node.Else.Body.Stmt[i+1:]...)
							i--
							continue
						}
						if st, ok := newChild.(ast.Statement); ok {
							node.Else.Body.Stmt[i] = st
						}
					}
					ctx.PopExprContext()
				} else {
					// Plain else body
					ctx.PushExprContext(ExprContext{Kind: CtxElseBody, ArgIdx: -1})
					if _, err := walkNode(node.Else, listener, adjuster, ctx); err != nil {
						return nil, err
					}
					ctx.PopExprContext()
				}
			}
		}
	case *stmt.ForEach:
		// foreach iterable expression
		ctx.PushExprContext(ExprContext{Kind: CtxForEachCond, ArgIdx: -1})
		_ = walkExpr(node.Set, listener, adjuster, ctx)
		ctx.PopExprContext()

		// Enter new scope and mark foreach item as local
		ctx.EnterScope()
		if node.Item != nil && node.Item.ID != "" {
			ctx.MarkLocal(node.Item.ID)
		}

		// foreach body
		ctx.PushExprContext(ExprContext{Kind: CtxForEachBody, ArgIdx: -1})
		for i, child := range node.Body.Stmt {
			newChild, err := walkNode(child, listener, adjuster, ctx)
			if err != nil {
				return nil, err
			}
			if newChild == nil {
				node.Body.Stmt = append(node.Body.Stmt[:i], node.Body.Stmt[i+1:]...)
				i--
				continue
			}
			if st, ok := newChild.(ast.Statement); ok {
				node.Body.Stmt[i] = st
			}
		}
		ctx.PopExprContext()
		ctx.ExitScope()
	case *stmt.ForLoop:
		if node.Init != nil {
			if initNode, ok := node.Init.(ast.Node); ok {
				ctx.PushExprContext(ExprContext{Kind: CtxForLoopInit, ArgIdx: -1})
				newInit, err := walkNode(initNode, listener, adjuster, ctx)
				ctx.PopExprContext()
				if err != nil {
					return nil, err
				}
				if newInit == nil {
					node.Init = nil
				} else if st, ok2 := newInit.(ast.Statement); ok2 {
					node.Init = st
				}
			}
		}
		if node.Cond != nil {
			ctx.PushExprContext(ExprContext{Kind: CtxForLoopCond, ArgIdx: -1})
			_ = walkExpr(node.Cond, listener, adjuster, ctx)
			ctx.PopExprContext()
		}
		if node.Post != nil {
			if postNode, ok := node.Post.(ast.Node); ok {
				ctx.PushExprContext(ExprContext{Kind: CtxForLoopPost, ArgIdx: -1})
				if _, err := walkNode(postNode, listener, adjuster, ctx); err != nil {
					ctx.PopExprContext()
					return nil, err
				}
				ctx.PopExprContext()
			}
		}
		for i, child := range node.Body.Stmt {
			newChild, err := walkNode(child, listener, adjuster, ctx)
			if err != nil {
				return nil, err
			}
			if st, ok := newChild.(ast.Statement); ok {
				node.Body.Stmt[i] = st
			}
		}
	case *stmt.Evaluate:
		handleEvaluate(node, listener, adjuster, ctx)
	case *stmt.Statement:
		if node.X != nil {
			ctx.PushExprContext(ExprContext{Kind: CtxSetLHS, ArgIdx: -1})
			_ = walkExpr(node.X, listener, adjuster, ctx)
			ctx.PopExprContext()
			// Record #set local variable on assignment when LHS is a selector name
			if sel, ok := node.X.(*aexpr.Select); ok {
				ctx.MarkLocal(sel.ID)
			}
		}
		if node.Y != nil {
			ctx.PushExprContext(ExprContext{Kind: CtxSetRHS, ArgIdx: -1})
			_ = walkExpr(node.Y, listener, adjuster, ctx)
			ctx.PopExprContext()
		}
	case *stmt.Append:
		// plain text context for this node
		ctx.PushExprContext(ExprContext{Kind: CtxPlainText, ArgIdx: -1})
		ctx.PopExprContext()
	}
	ctx.PopNode()
	exit := Event{Type: EventExitNode, Node: n, Context: ctx, ExprContext: ctx.CurrentExprContext()}
	if s, ok := ctx.GetSpan(n); ok {
		exit.Span = s
		exit.Position = ctx.ResolvePosition(s)
	}
	if listener != nil {
		listener.OnEvent(exit)
	}
	return n, nil
}

// walkExpr traverses expression AST and fires events/adjusters via synthetic nodes.
func walkExpr(e ast.Expression, listener ParserListener, adjuster NodeAdjuster, ctx *ParserContext) ast.Expression {
	if e == nil {
		return nil
	}
	// If no explicit context has been set, treat as append expression.
	if ctx.CurrentExprContext().Kind == CtxUnknown {
		ctx.PushExprContext(ExprContext{Kind: CtxAppendExpr, ArgIdx: -1})
		defer ctx.PopExprContext()
	}

	// Track per-name occurrence for selectors.
	occ := 0
	if sel, ok := e.(*aexpr.Select); ok {
		occ = ctx.BumpOccurrence(sel.ID)
	}

	evEnter := Event{Type: EventEnterNode, Node: e, Context: ctx, ExprContext: ctx.CurrentExprContext(), Occurrence: occ}
	if s, ok := ctx.GetSpan(e); ok {
		evEnter.Span = s
		evEnter.Position = ctx.ResolvePosition(s)
	}
	if listener != nil {
		listener.OnEvent(evEnter)
	}
	if adjuster != nil {
		act, err := adjuster.Adjust(e, ctx)
		if err == nil {
			switch act.Kind {
			case ActionReplace:
				if ex, ok := act.Node.(ast.Expression); ok {
					e = ex
				}
				if act.SkipChildren {
					ev := Event{Type: EventExitNode, Node: e, Context: ctx, ExprContext: ctx.CurrentExprContext(), Occurrence: occ}
					if s, ok := ctx.GetSpan(e); ok {
						ev.Span = s
						ev.Position = ctx.ResolvePosition(s)
					}
					if listener != nil {
						listener.OnEvent(ev)
					}
					return e
				}
			case ActionRemove:
				ev := Event{Type: EventExitNode, Node: e, Context: ctx, ExprContext: ctx.CurrentExprContext(), Occurrence: occ}
				if s, ok := ctx.GetSpan(e); ok {
					ev.Span = s
					ev.Position = ctx.ResolvePosition(s)
				}
				if listener != nil {
					listener.OnEvent(ev)
				}
				return nil
			case ActionPatchOnly:
				if act.SkipChildren {
					ev := Event{Type: EventExitNode, Node: e, Context: ctx, ExprContext: ctx.CurrentExprContext(), Occurrence: occ}
					if s, ok := ctx.GetSpan(e); ok {
						ev.Span = s
						ev.Position = ctx.ResolvePosition(s)
					}
					if listener != nil {
						listener.OnEvent(ev)
					}
					return e
				}
			case ActionKeep:
				// continue
			}
		}
	}
	switch x := e.(type) {
	case *aexpr.Select:
		if x.X != nil {
			_ = walkExpr(x.X, listener, adjuster, ctx)
		}
	case *aexpr.Call:
		for i := range x.Args {
			ctx.PushExprContext(ExprContext{Kind: CtxFuncArg, ArgIdx: int16(i)})
			_ = walkExpr(x.Args[i], listener, adjuster, ctx)
			ctx.PopExprContext()
		}
		if x.X != nil {
			_ = walkExpr(x.X, listener, adjuster, ctx)
		}
	case *aexpr.Binary:
		_ = walkExpr(x.X, listener, adjuster, ctx)
		_ = walkExpr(x.Y, listener, adjuster, ctx)
	case *aexpr.Unary:
		_ = walkExpr(x.X, listener, adjuster, ctx)
	case *aexpr.Parentheses:
		_ = walkExpr(x.P, listener, adjuster, ctx)
	case *aexpr.Range:
		// literals are terminal
	case *aexpr.Literal:
		// terminal
	}
	evExit := Event{Type: EventExitNode, Node: e, Context: ctx, ExprContext: ctx.CurrentExprContext(), Occurrence: occ}
	if s, ok := ctx.GetSpan(e); ok {
		evExit.Span = s
		evExit.Position = ctx.ResolvePosition(s)
	}
	if listener != nil {
		listener.OnEvent(evExit)
	}
	return e
}

// handleEvaluate applies Evaluate handling respecting context configuration.
func handleEvaluate(ev *stmt.Evaluate, listener ParserListener, adjuster NodeAdjuster, ctx *ParserContext) {
	// Always emit events for the Evaluate node expression itself
	ctx.PushExprContext(ExprContext{Kind: CtxEvaluate, ArgIdx: -1})
	_ = walkExpr(ev.X, listener, adjuster, ctx)
	ctx.PopExprContext()

	mode := ctx.EvalConfig.Mode
	if mode == EvalOpaque {
		return
	}
	// depth guard
	if ctx.EvalConfig.DepthLimit > 0 && ctx.EvalDepth >= ctx.EvalConfig.DepthLimit {
		// emit nothing further in inspect/rewrite when depth exceeded
		return
	}

	// Only allow string literal input unless safety permits selector whitelisting
	switch x := ev.X.(type) {
	case *aexpr.Literal:
		// ensure it's a string literal
		if x.RType == nil || x.RType.Kind() != reflect.String {
			return
		}
		data := []byte(x.Value)
		if max := ctx.EvalConfig.Safety.MaxBytes; max > 0 && len(data) > max {
			return
		}
		// Parse and traverse child template with spans to preserve positions
		block, spans, err := parser.ParseWithSpansDetailed(data)
		if err != nil {
			return
		}
		// Descend with incremented depth; use synthetic file label
		saved := ctx.EvalDepth
		ctx.EvalDepth = saved + 1
		childCtx := &ParserContext{Scope: NewScope(nil)}
		childCtx.InitSource(ctx.File+"#eval", data)
		childCtx.EvalDepth = ctx.EvalDepth
		_, _ = applyParserHooksWithSourceAndSpans(childCtx, block, listener, adjuster, spans)
		ctx.EvalDepth = saved
	case *aexpr.Select:
		// Only inspect if allowed by safety and whitelist (no rewrite possible without actual data)
		if ctx.EvalConfig.Safety.StringLiteralsOnly {
			return
		}
		id := x.ID
		if len(ctx.EvalConfig.Whitelist) > 0 && !contains(ctx.EvalConfig.Whitelist, id) {
			return
		}
		if contains(ctx.EvalConfig.Blacklist, id) {
			return
		}
		// In Inspect mode we do nothing more; in Rewrite mode as well, since we can't parse dynamic content here
		return
	default:
		// other dynamic expressions ignored
		return
	}
}

func reflectStringKind() int {
	// local helper to avoid importing reflect here explicitly
	// reflect.String has Kind() == 24; but safer is to use reflect in top scope; compromise: return 24
	return 24
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
