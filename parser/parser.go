package parser

import (
	"fmt"
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
)

func Parse(input []byte) (*stmt.Block, error) {
	if len(input) == 0 {
		return nil, nil
	}

	builder := NewBuilder()
	var tokenMatch *parsly.TokenMatch
	cursor := parsly.NewCursor("", input, 0)
	for cursor.Pos < len(input) {
		tokenMatch = cursor.MatchOne(SpecialSign)
		text := tokenMatch.Text(cursor)

		if tokenMatch.Code == parsly.EOF || cursor.Pos >= len(input) {
			if err := builder.PushStatement(appendToken, stmt.NewAppend(text)); err != nil {
				return nil, err
			}
			break
		}

		if cursor.Input[cursor.Pos] == '#' {
			cursor.MatchOne(NewLine)
			continue
		}

		err := appendStatementIfNeeded(text, builder)
		if err != nil {
			return nil, err
		}

		switch cursor.Input[cursor.Pos-1] {
		case '$':
			statement, err := matchSelector(cursor)
			if err != nil {
				return nil, err
			}
			builder.appendStatement(statement)

		case '#':
			statement, match, err := matchStatement(cursor)
			if err != nil {
				return nil, err
			}

			if err = builder.PushStatement(match, statement); err != nil {
				return nil, err
			}
		}
	}

	if builder.BufferSize() != 0 {
		return nil, fmt.Errorf("unterminated statements on the stack: %v", builder.buffer)
	}

	return builder.Block(), nil
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

func matchStatement(cursor *parsly.Cursor) (ast.Statement, int, error) {
	candidates := []*parsly.Token{If, ElseIf, Else, Set, ForEach, For, End}
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

		ifStmt, err := matchIf(expressionCursor)
		if err != nil {
			return nil, 0, err
		}
		return ifStmt, expressionCode, nil
	case elseToken:
		return &stmt.If{
			Condition: &expr.Binary{
				X:     expr.BoolExpression("true"),
				Token: "==",
				Y:     expr.BoolExpression("true"),
			},
			Body: stmt.Block{},
		}, expressionCode, nil

	case setToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		assignStmt, err := matchAssign(expressionCursor)
		if err != nil {
			return nil, expressionCode, err
		}

		return assignStmt, expressionCode, nil
	case forEachToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		forEachStmt, err := matchForEach(expressionCursor)
		if err != nil {
			return nil, 0, err
		}

		return forEachStmt, expressionCode, nil

	case forToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		forStmt, err := matchFor(expressionCursor)
		if err != nil {
			return nil, 0, err
		}

		return forStmt, expressionCode, nil

	case endToken:
		return nil, expressionCode, nil
	}

	return nil, 0, cursor.NewError(candidates...)
}
