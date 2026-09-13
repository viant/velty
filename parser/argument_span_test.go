package parser

import (
	"github.com/viant/velty/ast/expr"
	"strings"
	"testing"
)

func TestArgumentSelectorSpans(t *testing.T) {
	for _, source := range []string{`$dml.Insert("events", $unknown)`, `prefix $dml.Insert("events", $helper.Make($unknown))`, `prefix $rows[$unknown]`, `prefix $dml.Insert("events", $rows[$unknown])`} {
		_, spans, err := ParseWithSpansDetailed([]byte(source))
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for node, span := range spans {
			selector, ok := node.(*expr.Select)
			if !ok || selector.ID != "unknown" {
				continue
			}
			found = true
			start := strings.Index(source, "$unknown")
			if span.Start != start || span.End != start+len("$unknown")-1 {
				t.Fatalf("span %+v, expected %d..%d in %s", span, start, start+len("$unknown")-1, source)
			}
		}
		if !found {
			t.Fatal("selector span missing")
		}
	}
}

func TestFunctionCallSingleCharacterArguments(t *testing.T) {
	for _, tt := range []struct {
		source string
		count  int
	}{{`$F.Call(1)`, 1}, {`$F.Call(1, 2)`, 2}, {`$F.Call( )`, 0}} {
		root, err := Parse([]byte(tt.source))
		if err != nil {
			t.Fatal(err)
		}
		selector := root.Statements()[0].(*expr.Select)
		method := selector.X.(*expr.Select)
		call := method.X.(*expr.Call)
		if len(call.Args) != tt.count {
			t.Fatalf("%s args %d expected %d", tt.source, len(call.Args), tt.count)
		}
	}
}
