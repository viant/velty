package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
)

func matchVariable(cursor *parsly.Cursor) (*expr.Select, error) {
	candidates := []*parsly.Token{SelectorStart}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	if matched.Code == parsly.Invalid || matched.Code == parsly.EOF {
		return nil, cursor.NewError(candidates...)
	}

	variable, err := parseIdentity(cursor, false)
	if err != nil {
		return nil, err
	}
	return variable, nil
}

func matchSelector(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{SelectorBlock, Selector}
	matched := cursor.MatchAny(candidates...)
	if matched.Code == parsly.EOF {
		return nil, cursor.NewError(candidates...)
	}

	switch matched.Code {
	case selectorBlockToken:
		ID := matched.Text(cursor)

		selectorCursor := parsly.NewCursor("", []byte(ID[1:len(ID)-1]), 0)
		selector, err := parseIdentity(selectorCursor, true)
		if err != nil {
			return nil, err
		}
		return selector, nil

	case selectorToken:
		selectorValue := matched.Text(cursor)
		selectorCursor := parsly.NewCursor("", []byte(selectorValue), 0)
		selector, err := parseIdentity(selectorCursor, true)
		if err != nil {
			return nil, err
		}
		return selector, nil
	}

	return nil, cursor.NewError(candidates...)
}

func parseIdentity(cursor *parsly.Cursor, fullMatch bool) (*expr.Select, error) {
	var candidates []*parsly.Token
	if fullMatch {
		candidates = []*parsly.Token{ComplexSelector}
	} else {
		candidates = []*parsly.Token{NewVariable}
	}

	matched := cursor.MatchAny(candidates...)
	id := matched.Text(cursor)

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	}

	candidates = []*parsly.Token{Parentheses}
	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var call *expr.Call
	if matched.Code == parentheses {
		callValue := matched.Text(cursor)
		callCursor := parsly.NewCursor("", []byte(callValue[1:len(callValue)-1]), 0)

		var err error
		call, err = matchFunctionCall(callCursor)
		if err != nil {
			return nil, err
		}
	}

	return &expr.Select{
		ID:   id,
		Call: call,
	}, nil
}

func matchFunctionCall(cursor *parsly.Cursor) (*expr.Call, error) {
	expressions := make([]ast.Expression, 0)

	for cursor.Pos < cursor.InputSize-1 {
		argumentCursor := extractArgument(cursor)

		candidates := []*parsly.Token{Sub, Negation}
		matched := argumentCursor.MatchAfterOptional(WhiteSpace, candidates...)
		token := matchToken(matched)

		_, expression, err := matchOperand(argumentCursor, String, Boolean, Number)
		if err != nil {
			return nil, err
		}

		if token == "" {
			expressions = append(expressions, expression)
		} else {
			expressions = append(expressions, &expr.Unary{
				Token: token,
				X:     expression,
			})
		}

	}

	return &expr.Call{Args: expressions}, nil
}

func extractArgument(cursor *parsly.Cursor) *parsly.Cursor {
	candidates := []*parsly.Token{Coma}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case comaToken:
		return parsly.NewCursor("", cursor.Input[:cursor.Pos], 0)
	}

	return cursor
}
