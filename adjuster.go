package velty

import "github.com/viant/velty/ast"

// Patch represents a textual replacement for a span.
type Patch struct {
	Span        Span
	Replacement []byte
}

// ActionKind describes the result of an adjuster.
type ActionKind int

const (
	ActionKeep ActionKind = iota
	ActionReplace
	ActionRemove
	ActionPatchOnly
)

// Action is the result returned by an adjuster.
type Action struct {
	Kind         ActionKind
	Node         ast.Node // for replace
	Patches      []Patch  // optional patches
	SkipChildren bool     // skip walking children of this node
}

// Helpers
func Keep() Action              { return Action{Kind: ActionKeep} }
func Replace(n ast.Node) Action { return Action{Kind: ActionReplace, Node: n} }
func Remove() Action            { return Action{Kind: ActionRemove} }
func PatchSpan(s Span, repl []byte) Action {
	return Action{Kind: ActionPatchOnly, Patches: []Patch{{Span: s, Replacement: repl}}}
}
func (a Action) WithSkipChildren() Action { a.SkipChildren = true; return a }

// NodeAdjuster defines a transformation applied to AST nodes.
type NodeAdjuster interface {
	// Adjust applies transformations to the given node within the parser context.
	// It returns an action or error.
	Adjust(node ast.Node, ctx *ParserContext) (Action, error)
}

// AdjustFunc is a helper to define NodeAdjuster via a function.
type AdjustFunc func(node ast.Node, ctx *ParserContext) (Action, error)

// Adjust implements NodeAdjuster for AdjustFunc.
func (f AdjustFunc) Adjust(node ast.Node, ctx *ParserContext) (Action, error) {
	return f(node, ctx)
}

// AdjusterChain applies multiple NodeAdjusters in sequence.
type AdjusterChain struct {
	Adjusters []NodeAdjuster
}

// NewAdjusterChain constructs a chain of NodeAdjusters.
func NewAdjusterChain(adjusters ...NodeAdjuster) *AdjusterChain {
	return &AdjusterChain{Adjusters: adjusters}
}

// Adjust runs each adjuster on the node in order, passing the parser context.
// It merges patches and respects removal and replace.
func (c *AdjusterChain) Adjust(node ast.Node, ctx *ParserContext) (Action, error) {
	result := Keep()
	curr := node
	for _, adj := range c.Adjusters {
		act, err := adj.Adjust(curr, ctx)
		if err != nil {
			return Action{}, err
		}
		// merge patches
		if len(act.Patches) > 0 {
			ctx.AddPatches(act.Patches...)
		}
		// respect removal
		if act.Kind == ActionRemove {
			return act, nil
		}
		// handle replace
		if act.Kind == ActionReplace {
			curr = act.Node
			result = act
		}
		// propagate SkipChildren if any
		if act.SkipChildren {
			result.SkipChildren = true
		}
	}
	if result.Kind == ActionReplace {
		return result, nil
	}
	return Keep(), nil
}
