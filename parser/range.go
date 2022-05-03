package parser

import (
	"github.com/viant/parsly"
	ast2 "github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
)

func matchForEach(cursor *parsly.Cursor) (*stmt.ForEach, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}
	candidates := []*parsly.Token{ComaTerminator}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var index *expr.Select
	if matched.Code == comaToken {
		index, err = matchVariable(cursor)
		if err != nil {
			return nil, err
		}
	}
	candidates = []*parsly.Token{In}

	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)
	if matched.Code == parsly.Invalid || matched.Code == parsly.EOF {
		return nil, cursor.NewError(candidates...)
	}

	dataSet, err := matchRangeable(cursor)
	if err != nil {
		return nil, err
	}

	return &stmt.ForEach{
		Index: index,
		Item:  variable,
		Set:   dataSet,
		Body:  stmt.Block{},
	}, nil
}

func matchFor(cursor *parsly.Cursor) (*stmt.ForLoop, error) {
	initCursor := extractForSegment(cursor)
	forInit, err := matchAssign(initCursor)
	if err != nil {
		return nil, err
	}

	conditionCursor := extractForSegment(cursor)
	forCondition, err := matchEquationExpression(conditionCursor)
	if err != nil {
		return nil, err
	}

	forPostCursor := extractForSegment(cursor)
	forPost, err := matchForPost(forPostCursor)
	if err != nil {
		return nil, err
	}

	return &stmt.ForLoop{
		Init: forInit,
		Cond: forCondition,
		Post: forPost,
	}, nil
}

func matchForPost(cursor *parsly.Cursor) (ast2.Statement, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}

	candidates := []*parsly.Token{Increment, Decrement, AddEqual, Add, SubEqual, Sub, MultiplyEqual, Multiply, QuoEqual, Quo}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(candidates...)

	case decrementToken:
		return &stmt.Statement{
			X:  variable,
			Op: ast2.ASSIGN,
			Y:  expr.BinaryExpression(variable, ast2.SUB, expr.NumberLiteral("1")),
		}, nil

	case incrementToken:
		return &stmt.Statement{
			X:  variable,
			Op: ast2.ASSIGN,
			Y:  expr.BinaryExpression(variable, ast2.ADD, expr.NumberLiteral("1")),
		}, nil
	}

	token := matchToken(matched)
	_, rightOperand, err := matchOperand(cursor)
	if err != nil {
		return nil, err
	}

	token, rightOperand = normalizeTokensIfNeeded(variable, token, rightOperand, matched)

	return &stmt.Statement{
		X:  variable,
		Op: token,
		Y:  rightOperand,
	}, nil
}

func extractForSegment(cursor *parsly.Cursor) *parsly.Cursor {
	candidates := []*parsly.Token{ExpressionBlock, ExpressionStart, ExpressionEnd}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case expressionBlockToken:
		segment := matched.Text(cursor)
		return parsly.NewCursor("", []byte(segment[1:len(segment)-1]), 0)
	case expressionStartToken:
		expressionStartMatch := cursor.MatchAfterOptional(WhiteSpace, ExpressionEnd)
		if expressionStartMatch.Code == parsly.EOF || expressionStartMatch.Code == parsly.Invalid {
			return parsly.NewCursor("", cursor.Input[cursor.Pos:], 0)
		}

		text := expressionStartMatch.Text(cursor)
		return parsly.NewCursor("", []byte(text[:len(text)-1]), 0)
	case expressionEndToken:
		segment := matched.Text(cursor)
		return parsly.NewCursor("", []byte(segment), 0)
	}

	return parsly.NewCursor("", cursor.Input[cursor.Pos:], 0)
}
