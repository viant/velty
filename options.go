package velty

import (
	"github.com/viant/velty/functions"
)

// Option represents Planner generic option
type Option interface{}

// BufferSize represents initial size of the buffer
type BufferSize int

// CacheSize represents cache size in case of the dynamic template evaluation
type CacheSize int

// EscapeHTML escapes HTML in passed variables.
type EscapeHTML bool

// PanicOnError panics and recover when first error returned.
type PanicOnError bool

// TypeParser parses type string representation into reflect.Type
type TypeParser = functions.TypeParser

// Listener registers a ParserListener to receive parse events.
// Usage: New(Listener(yourListener))
type Listener ParserListener

// Adjuster registers a NodeAdjuster to transform AST nodes.
// Usage: New(Adjuster(yourAdjuster))
type Adjuster NodeAdjuster

// Policies registers a PolicyRegistry to drive adjustments.
// Usage: New(Policies(reg)) and then use reg.AsAdjuster() as Adjuster.
type Policies *PolicyRegistry

// Evaluate controls evaluate traversal.
// EvalCfg registers evaluate traversal configuration.
type EvalCfg EvaluateConfig

// PlanHooks registers a planning listener to receive binding hooks.
// Usage: New(PlanHooks(yourPlannerListener))
type PlanHooks PlannerListener
