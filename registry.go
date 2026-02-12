package velty

import (
	"github.com/viant/velty/ast"
	"sort"
)

// PolicyRegistry stores and applies policies in priority order.
type PolicyRegistry struct {
	items []Policy
	idx   map[string]int
}

func NewPolicyRegistry() *PolicyRegistry {
	return &PolicyRegistry{items: make([]Policy, 0, 8), idx: map[string]int{}}
}

// Register adds or replaces a policy by name.
func (r *PolicyRegistry) Register(p Policy) {
	if p == nil {
		return
	}
	if i, ok := r.idx[p.Name()]; ok {
		r.items[i] = p
		r.sort()
		return
	}
	r.items = append(r.items, p)
	r.idx[p.Name()] = len(r.items) - 1
	r.sort()
}

// Enable toggles a policy by name when it implements BasicPolicy convention.
func (r *PolicyRegistry) Enable(name string, on bool) {
	if i, ok := r.idx[name]; ok {
		// best effort: attempt to toggle Active if policy is *BasicPolicy
		if bp, ok2 := r.items[i].(*BasicPolicy); ok2 {
			bp.Active = on
		}
	}
}

func (r *PolicyRegistry) Policies() []Policy { return r.items }

func (r *PolicyRegistry) sort() {
	sort.SliceStable(r.items, func(i, j int) bool { return r.items[i].Priority() < r.items[j].Priority() })
	// rebuild index
	for i := range r.items {
		r.idx[r.items[i].Name()] = i
	}
}

// AsAdjuster builds a NodeAdjuster that evaluates registered policies.
func (r *PolicyRegistry) AsAdjuster() NodeAdjuster {
	return AdjustFunc(func(node ast.Node, ctx *ParserContext) (Action, error) {
		var replaced ast.Node = node
		skipChildren := false
		for _, p := range r.items {
			if !p.Enabled() {
				continue
			}
			if !p.Match(replaced, ctx) {
				continue
			}
			act, err := p.Apply(replaced, ctx)
			if err != nil {
				return Action{}, err
			}
			if len(act.Patches) > 0 {
				ctx.AddPatches(act.Patches...)
			}
			switch act.Kind {
			case ActionKeep:
				// continue to next policy
			case ActionRemove:
				return act, nil
			case ActionReplace:
				if act.Node != nil {
					replaced = act.Node
				}
				if act.SkipChildren {
					skipChildren = true
				}
			case ActionPatchOnly:
				if act.SkipChildren {
					skipChildren = true
				}
			}
		}
		if replaced != node || skipChildren {
			res := Replace(replaced)
			res.SkipChildren = skipChildren
			return res, nil
		}
		return Keep(), nil
	})
}
