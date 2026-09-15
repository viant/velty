package parser

import (
	"testing"

	"github.com/viant/parsly"
)

func TestMatchStatementHostBoundary(t *testing.T) {
	for _, template := range []string{
		`#foreach($id in $IDs)$id#if($foreach.HasNext),#end#end`,
		`#if($ok)$id#else 0#end`,
		`#foreach($group in $Groups)#foreach($id in $group.IDs)$id#end#end`,
		`#set($value = 1)`,
		`#if($ok)literal #unknown(1) and $1#end`,
	} {
		t.Run(template, func(t *testing.T) {
			input := "prefix " + template + ", 42) ORDER BY ID"
			cursor := parsly.NewCursor("", []byte(input), 0)
			cursor.Pos = len("prefix ")
			actual, err := MatchStatement(cursor)
			if err != nil || actual == nil {
				t.Fatalf("statement=%T err=%v", actual, err)
			}
			if cursor.Pos != len("prefix ")+len(template) {
				t.Fatalf("consumed %q", input[:cursor.Pos])
			}
			if _, err := Parse([]byte(template)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMatchStatementRejectsMalformed(t *testing.T) {
	for _, input := range []string{
		`#`, `#[[literal]]#`, `text`, `#unknown($id)`, `#end`, `#else`, `#foreach($id in $IDs)$id`,
		`#foreach($id $IDs)$id#end`, `#if($ok)$id`,
		`#foreach($id in $IDs)#if()$id#end#end`,
	} {
		t.Run(input, func(t *testing.T) {
			cursor := parsly.NewCursor("", []byte("prefix "+input), 0)
			cursor.Pos = len("prefix ")
			if _, err := MatchStatement(cursor); err == nil {
				t.Fatal("expected error")
			}
			if cursor.Pos != len("prefix ") {
				t.Fatal("failed parse advanced cursor")
			}
		})
	}
}

func TestParseRetainsLiteralUnknownDirective(t *testing.T) {
	root, err := Parse([]byte(`before #unknown(1) after`))
	if err != nil || root == nil {
		t.Fatal(err)
	}
	if len(root.Stmt) == 0 {
		t.Fatal("literal text lost")
	}
}
