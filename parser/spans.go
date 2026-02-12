package parser

import (
	"github.com/viant/velty/ast"
)

// NodeSpan represents a byte range [Start, End] in the parsed source.
type NodeSpan struct {
	Start int
	End   int
}

// nodeSpans holds spans for the most recent parse when span recording is enabled.
var nodeSpans map[ast.Node]NodeSpan

// recordSpans controls whether spans are recorded for the current parse.
var recordSpans bool

// resetSpans initializes the global span registry for a new parse when recording is enabled.
func resetSpans() {
	if !recordSpans {
		nodeSpans = nil
		return
	}
	nodeSpans = make(map[ast.Node]NodeSpan)
}

// recordSpan registers a span for the given node when recording is enabled.
func recordSpan(n ast.Node, start, end int) {
	if !recordSpans || n == nil {
		return
	}
	if nodeSpans == nil {
		nodeSpans = make(map[ast.Node]NodeSpan)
	}
	nodeSpans[n] = NodeSpan{Start: start, End: end}
}

// getSpan returns a span for a node if recorded.
func getSpan(n ast.Node) (NodeSpan, bool) {
	if nodeSpans == nil {
		return NodeSpan{}, false
	}
	s, ok := nodeSpans[n]
	return s, ok
}

// Spans exposes the current parse node spans.
func Spans() map[ast.Node]NodeSpan {
	return nodeSpans
}
