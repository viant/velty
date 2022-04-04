package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/internal/ast"
	aexpr "github.com/viant/velty/internal/ast/expr"
	astmt "github.com/viant/velty/internal/ast/stmt"
)

func matchForEach(cursor *parsly.Cursor) (*astmt.ForEach, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}
	candidates := []*parsly.Token{ComaTerminator}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var index *aexpr.Select
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

	return &astmt.ForEach{
		Index: index,
		Item:  variable,
		Set:   dataSet,
		Body:  astmt.Block{},
	}, nil
}

func matchFor(cursor *parsly.Cursor) (*astmt.ForLoop, error) {
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

	return &astmt.ForLoop{
		Init: forInit,
		Cond: forCondition,
		Post: forPost,
	}, nil
}

func matchForPost(cursor *parsly.Cursor) (ast.Statement, error) {
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
		return &astmt.Statement{
			X:  variable,
			Op: ast.ASSIGN,
			Y:  aexpr.BinaryExpression(variable, ast.SUB, aexpr.NumberLiteral("1")),
		}, nil

	case incrementToken:
		return &astmt.Statement{
			X:  variable,
			Op: ast.ASSIGN,
			Y:  aexpr.BinaryExpression(variable, ast.ADD, aexpr.NumberLiteral("1")),
		}, nil
	}

	token := matchToken(matched)
	_, rightOperand, err := matchOperand(cursor)
	if err != nil {
		return nil, err
	}

	token, rightOperand = normalizeTokensIfNeeded(variable, token, rightOperand, matched)

	return &astmt.Statement{
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
