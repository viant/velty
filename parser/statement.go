package parser

import (
	"fmt"

	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/stmt"
)

// MatchStatement consumes one directive, including its nested block, at the
// cursor's current '#' byte. It leaves following host-language text untouched.
// Malformed recognized directives are errors; literal body text follows Parse.
// The cursor is unchanged on failure.
func MatchStatement(cursor *parsly.Cursor) (ast.Statement, error) {
	if cursor == nil || cursor.Pos < 0 || cursor.Pos >= len(cursor.Input) || cursor.Input[cursor.Pos] != '#' {
		return nil, fmt.Errorf("expected template directive")
	}
	local := *cursor
	block, err := parseCursor(&local, nil, true)
	if err != nil {
		return nil, err
	}
	if len(block.Stmt) != 1 {
		return nil, fmt.Errorf("expected one template statement")
	}
	if _, literal := block.Stmt[0].(*stmt.Append); literal {
		return nil, fmt.Errorf("expected template directive, got literal text")
	}
	cursor.Pos = local.Pos
	return block.Stmt[0], nil
}
