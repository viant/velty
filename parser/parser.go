package parser

import (
	"fmt"
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
	"strings"
)

// parse is the internal implementation shared by Parse and ParseWithSpans*.
func parse(input []byte, spans *spanState) (*stmt.Block, error) {
	if len(input) == 0 {
		return &stmt.Block{}, nil
	}

	builder := NewBuilder()
	var tokenMatch *parsly.TokenMatch
	cursor := parsly.NewCursor("", input, 0)
outer:
	for cursor.Pos < len(input) {
		tokenMatch = cursor.MatchOne(SpecialSign)
		text := tokenMatch.Text(cursor)

		if tokenMatch.Code == parsly.EOF || cursor.Pos >= len(input) {
			if err := builder.PushStatement(appendToken, stmt.NewAppend(text)); err != nil {
				return nil, cursorErr(cursor, err)
			}
			break
		}

		if cursor.Input[cursor.Pos] == '#' {
			cursor.MatchOne(NewLine)
			continue
		}

		err := appendStatementIfNeeded(text, builder)
		if err != nil {
			return nil, cursorErr(cursor, err)
		}

		lastPosition := cursor.Pos - 1
		switch cursor.Input[cursor.Pos-1] {
		case '$':
			statement, err := MatchSelector(cursor, spans)
			if err != nil {
				rawValue := cursor.Input[lastPosition:cursor.Pos]
				if errr := builder.PushStatement(appendToken, stmt.NewAppend(string(rawValue))); errr != nil {
					return nil, cursorErr(cursor, errr)
				}
				continue outer
			}
			builder.appendStatement(statement)

		case '#':
			appendStmt, ok := checkIfEscaped(cursor)
			if ok {
				if err = builder.PushStatement(appendToken, appendStmt); err != nil {
					return nil, err
				}
				continue
			}

			statement, match, err := matchStatement(cursor, spans)
			if err != nil {
				rawValue := cursor.Input[lastPosition:cursor.Pos]
				if errr := builder.PushStatement(appendToken, stmt.NewAppend(string(rawValue))); errr != nil {
					return nil, cursorErr(cursor, errr)
				}
				continue outer
			}

			if err = builder.PushStatement(match, statement); err != nil {
				return nil, cursorErr(cursor, err)
			}
		}
	}

	if builder.BufferSize() != 0 {
		return nil, fmt.Errorf("unterminated statements on the stack: %v", builder.buffer)
	}

	return builder.Block(), nil
}

// Parse parses the input template without recording spans (fast path).
func Parse(input []byte) (*stmt.Block, error) {
	return parse(input, nil)
}

// ParseWithSpans parses the input template and records node spans.
func ParseWithSpans(input []byte) (*stmt.Block, error) {
	root, _, err := ParseWithSpansDetailed(input)
	return root, err
}

// ParseWithSpansDetailed parses the input template and returns node spans.
func ParseWithSpansDetailed(input []byte) (*stmt.Block, map[ast.Node]NodeSpan, error) {
	spans := newSpanState(true)
	root, err := parse(input, spans)
	if err != nil {
		return nil, nil, err
	}
	return root, spans.Spans(), nil
}

func checkIfEscaped(cursor *parsly.Cursor) (*stmt.Append, bool) {
	lastCursorPos := cursor.Pos
	matched := cursor.MatchOne(SquareBrackets)
	if matched.Code == squareBracketsToken {
		body := matched.Text(cursor)
		if strings.HasPrefix(body, "[[") && strings.HasSuffix(body, "]]") {
			matched = cursor.MatchOne(Hash)
			if matched.Code == hashToken {
				return stmt.NewAppend(body[2 : len(body)-2]), true
			}
		}
	}

	cursor.Pos = lastCursorPos
	return nil, false
}

func cursorErr(cursor *parsly.Cursor, err error) error {
	return fmt.Errorf("%w, cursor position: %v", err, cursor.Pos)
}

func appendStatementIfNeeded(text string, stack *Builder) error {
	text = text[:len(text)-1]
	if len(text) == 0 {
		return nil
	}

	if err := stack.PushStatement(appendToken, stmt.NewAppend(text)); err != nil {
		return err
	}
	return nil
}

func matchStatement(cursor *parsly.Cursor, spans *spanState) (ast.Statement, int, error) {
	matched := cursor.MatchAfterOptional(WhiteSpace, Brackets)
	if matched.Token.Code == bracketsToken {
		stmt := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(stmt[1:len(stmt)-1]), 0)
		return matchStatement(newCursor, spans)
	}

	candidates := []*parsly.Token{If, ElseIf, Else, Set, ForEach, For, Break, Evaluate, End}
	expressionMatch := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	expressionCode := expressionMatch.Code

	switch expressionMatch.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, 0, cursor.NewError(candidates...)
	case ifToken, elseIfToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		ifStmt, err := matchIf(expressionCursor, spans)
		if err != nil {
			return nil, 0, err
		}
		return ifStmt, expressionCode, nil
	case elseToken:
		return &stmt.If{
			Condition: &expr.Binary{
				X:     expr.BoolLiteral("true"),
				Token: "==",
				Y:     expr.BoolLiteral("true"),
			},
			Body: stmt.Block{},
		}, expressionCode, nil

	case setToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		assignStmt, err := matchAssign(expressionCursor, spans)
		if err != nil {
			return nil, expressionCode, err
		}

		return assignStmt, expressionCode, nil
	case forEachToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		forEachStmt, err := matchForEach(expressionCursor, spans)
		if err != nil {
			return nil, 0, err
		}

		return forEachStmt, expressionCode, nil

	case forToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		forStmt, err := matchFor(expressionCursor, spans)
		if err != nil {
			return nil, 0, err
		}

		return forStmt, expressionCode, nil

	case breakToken:
		return &stmt.Break{}, expressionCode, nil

	case evaluateToken:
		evaluateCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}
		_, operand, err := matchOperand(evaluateCursor, spans, String)

		if err != nil {
			return nil, 0, err
		}

		return &stmt.Evaluate{X: operand}, expressionCode, nil
	case endToken:
		return nil, expressionCode, nil
	}

	return nil, 0, cursor.NewError(candidates...)
}
