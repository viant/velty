package velty

import "github.com/viant/velty/ast"

// EventType represents the type of parser event.
type EventType int

const (
	// EventEnterNode indicates entering an AST node.
	EventEnterNode EventType = iota
	// EventExitNode indicates exiting an AST node.
	EventExitNode
)

// Event carries parser event data to listeners.
type Event struct {
	Type    EventType
	Node    ast.Node
	Context *ParserContext
	// Optional span/position if available
	Span     Span
	Position Position
	// ExprContext describes where this node appears (control/text/etc.).
	ExprContext ExprContext
	// Occurrence is a 1-based occurrence index for variable-like symbols
	// such as selectors ($X) when applicable; otherwise 0.
	Occurrence int
}

// ParserListener defines hooks for parser events during AST construction.
// Implementations can perform actions when nodes are entered or exited.
type ParserListener interface {
	// OnEvent is invoked for each parser event with full event data.
	OnEvent(e Event)
}
