package velty

import (
	"testing"

	"github.com/stretchr/testify/require"
	aexpr "github.com/viant/velty/ast/expr"
	"github.com/viant/velty/parser"
)

type ctxCaptureListener struct{ events []Event }

func (c *ctxCaptureListener) OnEvent(e Event) { c.events = append(c.events, e) }

func TestExprContext_BasicPatterns(t *testing.T) {
	t.Run("if_condition_selector", func(t *testing.T) {
		tpl := []byte("#if($X)ok#end")
		cap := &ctxCaptureListener{}
		root, err := parser.ParseWithSpans(tpl)
		require.NoError(t, err)
		_, _, err = applyParserHooksWithConfig("", tpl, root, cap, nil, EvaluateConfig{})
		require.NoError(t, err)

		var selCtx ExprContextKind
		for _, ev := range cap.events {
			if ev.Type != EventEnterNode {
				continue
			}
			if _, ok := ev.Node.(*aexpr.Select); ok {
				selCtx = ev.ExprContext.Kind
				break
			}
		}
		require.Equal(t, CtxIfCond, selCtx, "selector in #if condition should be CtxIfCond")
	})

	t.Run("foreach_iterable_and_body", func(t *testing.T) {
		tpl := []byte("#foreach($x in $items)$x#end")
		cap := &ctxCaptureListener{}
		root, err := parser.ParseWithSpans(tpl)
		require.NoError(t, err)
		_, _, err = applyParserHooksWithConfig("", tpl, root, cap, nil, EvaluateConfig{})
		require.NoError(t, err)

		var iterCtx, bodyCtx ExprContextKind
		foundIter := false
		foundBody := false
		for _, ev := range cap.events {
			if ev.Type != EventEnterNode {
				continue
			}
			if sel, ok := ev.Node.(*aexpr.Select); ok {
				switch sel.ID {
				case "items":
					iterCtx = ev.ExprContext.Kind
					foundIter = true
				case "x":
					bodyCtx = ev.ExprContext.Kind
					foundBody = true
				}
			}
		}
		require.True(t, foundIter, "foreach iterable selector not found")
		require.True(t, foundBody, "foreach body selector not found")
		require.Equal(t, CtxForEachCond, iterCtx, "foreach set expression should be CtxForEachCond")
		require.Equal(t, CtxForEachBody, bodyCtx, "foreach body selector should be CtxForEachBody")
	})

	t.Run("set_lhs_rhs", func(t *testing.T) {
		tpl := []byte("#set($x = $y)")
		cap := &ctxCaptureListener{}
		root, err := parser.ParseWithSpans(tpl)
		require.NoError(t, err)
		_, _, err = applyParserHooksWithConfig("", tpl, root, cap, nil, EvaluateConfig{})
		require.NoError(t, err)

		var lhsCtx, rhsCtx ExprContextKind
		for _, ev := range cap.events {
			if ev.Type != EventEnterNode {
				continue
			}
			sel, ok := ev.Node.(*aexpr.Select)
			if !ok {
				continue
			}
			switch sel.ID {
			case "x":
				lhsCtx = ev.ExprContext.Kind
			case "y":
				rhsCtx = ev.ExprContext.Kind
			}
		}
		require.Equal(t, CtxSetLHS, lhsCtx, "set LHS should be CtxSetLHS")
		require.Equal(t, CtxSetRHS, rhsCtx, "set RHS should be CtxSetRHS")
	})

	t.Run("call_arguments", func(t *testing.T) {
		// TODO: add a stable call-arg context test once
		// function/method call shapes are finalized for instrumentation.
	})

	t.Run("evaluate_expression_context", func(t *testing.T) {
		tpl := []byte("#evaluate(\"$X\")")
		cap := &ctxCaptureListener{}
		root, err := parser.ParseWithSpans(tpl)
		require.NoError(t, err)
		evalCfg := EvaluateConfig{Mode: EvalInspect}
		_, _, err = applyParserHooksWithConfig("", tpl, root, cap, nil, evalCfg)
		require.NoError(t, err)

		var evalCtx ExprContextKind
		for _, ev := range cap.events {
			if ev.Type != EventEnterNode {
				continue
			}
			if _, ok := ev.Node.(*aexpr.Select); ok {
				evalCtx = ev.ExprContext.Kind
				break
			}
		}
		require.Equal(t, CtxEvaluate, evalCtx, "selector inside #evaluate should be CtxEvaluate")
	})

	// TODO: add explicit elseif/else and for-loop
	// context assertions once the parser representation
	// and walker semantics are fully locked in.

	t.Run("varkind_and_occurrence", func(t *testing.T) {
		tpl := []byte("$P $P #set($l = $P) $l")
		cap := &ctxCaptureListener{}
		root, err := parser.ParseWithSpans(tpl)
		require.NoError(t, err)
		seeds := &SymbolSeeds{Params: []string{"P"}}
		_, _, err = applyParserHooksWithConfig("", tpl, root, cap, nil, EvaluateConfig{}, seeds)
		require.NoError(t, err)

		counts := map[string][]Event{}
		for _, ev := range cap.events {
			if ev.Type != EventEnterNode {
				continue
			}
			sel, ok := ev.Node.(*aexpr.Select)
			if !ok {
				continue
			}
			counts[sel.ID] = append(counts[sel.ID], ev)
		}

		// P should be classified as VarParam
		pEvents := counts["P"]
		require.GreaterOrEqual(t, len(pEvents), 3)
		for i, ev := range pEvents {
			kind := ev.Context.VarKind("P")
			require.Equal(t, VarParam, kind, "P should be VarParam")
			_ = i // occurrence is best-effort; not asserted here
		}

		// l should be a local (from #set)
		lEvents := counts["l"]
		require.GreaterOrEqual(t, len(lEvents), 1)
		for i, ev := range lEvents {
			kind := ev.Context.VarKind("l")
			require.Equal(t, VarLocal, kind, "l should be VarLocal")
			_ = i
		}
	})
}
