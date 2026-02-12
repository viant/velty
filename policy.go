package velty

import (
	"github.com/viant/velty/ast"
	"reflect"
)

// PolicyKind is a string label for node kinds (e.g., "Selector", "Call", "If", "ForEach", ...)
type PolicyKind string

// Policy represents a composable rule applied over nodes with metadata.
type Policy interface {
	// Name is the policy identifier.
	Name() string
	// Priority defines policy ordering (lower first).
	Priority() int
	// Enabled reports whether policy is active.
	Enabled() bool
	// Match selects nodes/contexts this policy applies to.
	Match(node ast.Node, ctx *ParserContext) bool
	// Apply executes policy logic and returns an action.
	Apply(node ast.Node, ctx *ParserContext) (Action, error)
}

// NodeKindOf returns a coarse node kind name for policies.
func NodeKindOf(n ast.Node) PolicyKind {
	if n == nil {
		return ""
	}
	t := reflect.TypeOf(n)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return PolicyKind(t.Name())
}

// BasicPolicy is a helper policy base with simple filters and a function logic.
type BasicPolicy struct {
	ID       string
	Order    int
	Active   bool
	Kinds    []PolicyKind // optional node kinds filter
	Contexts []string     // optional context kinds (reserved)
	Fn       func(node ast.Node, ctx *ParserContext) (Action, error)
}

func (p *BasicPolicy) Name() string  { return p.ID }
func (p *BasicPolicy) Priority() int { return p.Order }
func (p *BasicPolicy) Enabled() bool { return p.Active }

func (p *BasicPolicy) Match(node ast.Node, ctx *ParserContext) bool {
	if !p.Active {
		return false
	}
	if len(p.Kinds) == 0 {
		return true
	}
	kind := NodeKindOf(node)
	for _, k := range p.Kinds {
		if k == kind {
			return true
		}
	}
	return false
}

func (p *BasicPolicy) Apply(node ast.Node, ctx *ParserContext) (Action, error) {
	if p.Fn == nil {
		return Keep(), nil
	}
	return p.Fn(node, ctx)
}
