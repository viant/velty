package parser

import (
	"github.com/viant/velty/ast"
)

// NodeSpan represents a byte range [Start, End] in the parsed source.
type NodeSpan struct {
	Start int
	End   int
}

type spanState struct {
	record bool
	spans  map[ast.Node]NodeSpan
}

func newSpanState(record bool) *spanState {
	st := &spanState{record: record}
	if record {
		st.spans = make(map[ast.Node]NodeSpan)
	}
	return st
}

// recordSpan registers a span for the given node when recording is enabled.
func (s *spanState) recordSpan(n ast.Node, start, end int) {
	if s == nil || !s.record || n == nil {
		return
	}
	if s.spans == nil {
		s.spans = make(map[ast.Node]NodeSpan)
	}
	s.spans[n] = NodeSpan{Start: start, End: end}
}

// getSpan returns a span for a node if recorded.
func (s *spanState) getSpan(n ast.Node) (NodeSpan, bool) {
	if s == nil || s.spans == nil {
		return NodeSpan{}, false
	}
	span, ok := s.spans[n]
	return span, ok
}

// Spans exposes parse node spans for this parse invocation.
func (s *spanState) Spans() map[ast.Node]NodeSpan {
	if s == nil || s.spans == nil {
		return nil
	}
	return s.spans
}
