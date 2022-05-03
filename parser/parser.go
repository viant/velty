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
			statement, err := matchSelector(cursor)
			if err != nil {
				rawValue := cursor.Input[lastPosition:cursor.Pos]
				if errr := builder.PushStatement(appendToken, stmt.NewAppend(string(rawValue))); errr != nil {
					return nil, cursorErr(cursor, errr)
				}
				continue outer
			}
			builder.appendStatement(statement)

		case '#':
			statement, match, err := matchStatement(cursor)
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

func matchStatement(cursor *parsly.Cursor) (ast.Statement, int, error) {
	matched := cursor.MatchAfterOptional(WhiteSpace, Brackets)
	if matched.Token.Code == bracketsToken {
		stmt := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(stmt[1:len(stmt)-1]), 0)
		return matchStatement(newCursor)
	}

	candidates := []*parsly.Token{If, ElseIf, Else, Set, ForEach, For, Evaluate, End}
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

	case evaluateToken:
		evaluateCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}
		_, operand, err := matchOperand(evaluateCursor, String)

		if err != nil {
			return nil, 0, err
		}

		return &stmt.Evaluate{X: operand}, expressionCode, nil
	case endToken:
		return nil, expressionCode, nil
	}

	return nil, 0, cursor.NewError(candidates...)
}
