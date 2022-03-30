package parser

import (
	"fmt"
	"github.com/viant/parsly"
	"github.com/viant/velty/internal/ast"
	aexpr "github.com/viant/velty/internal/ast/expr"
	"github.com/viant/velty/internal/utils"
	"strconv"
	"strings"
)

func matchVariable(cursor *parsly.Cursor) (*aexpr.Select, error) {
	candidates := []*parsly.Token{SelectorStart}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case selectorStartToken:
		return parseIdentity(cursor)
	}
	return nil, cursor.NewError(candidates...)
}

func matchRangeable(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{SelectorStart, SquareBrackets}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case selectorStartToken:
		return parseIdentity(cursor)
	case squareBracketsToken:
		text := matched.Text(cursor)
		rangeCursor := parsly.NewCursor("", []byte(text[1:len(text)-1]), 0)

		return matchRange(rangeCursor)
	}
	return nil, cursor.NewError(candidates...)
}

func matchRange(cursor *parsly.Cursor) (ast.Expression, error) {
	variables := strings.Split(string(cursor.Input), "...")
	if len(variables) != 2 {
		return nil, fmt.Errorf("range expected to have two number literals but got %v", variables)
	}

	begin, err := strconv.Atoi(variables[0])
	if err != nil {
		return nil, err
	}

	finish, err := strconv.Atoi(variables[1])
	if err != nil {
		return nil, err
	}

	return &aexpr.Range{
		X: aexpr.Number(begin),
		Y: aexpr.Number(finish),
	}, nil
}

func matchSelector(cursor *parsly.Cursor) (ast.Expression, error) {
	matched := cursor.MatchOne(Negation) // Java velocity supports `$!`. If String is null, then it will be replaced with Empty String. In Go - we don't need that
	matched = cursor.MatchOne(SelectorBlock)

	if matched.Code == selectorBlockToken {
		ID := matched.Text(cursor)
		selectorCursor := parsly.NewCursor("", []byte(ID[1:len(ID)-1]), 0)
		result, err := matchSelector(selectorCursor)
		if err != nil {
			return nil, err
		}

		if selectorCursor.Pos < selectorCursor.InputSize {
			return nil, fmt.Errorf("expected to match all data, but couldn't match %v", string(cursor.Input[cursor.Pos:]))
		}

		return result, err

	}

	candidates := []*parsly.Token{SelectorBlock, Selector}
	matched = cursor.MatchAny(candidates...)
	if matched.Code == parsly.EOF {
		return nil, cursor.NewError(candidates...)
	}

	switch matched.Code {
	case selectorToken:
		selectorValue := matched.Text(cursor)
		selectorCursor := parsly.NewCursor("", []byte(selectorValue), 0)
		selector, err := parseIdentity(selectorCursor)
		if err != nil {
			return nil, err
		}

		selector.X, err = matchCall(cursor)
		if err != nil {
			return nil, err
		}
		return selector, nil
	}

	return nil, cursor.NewError(candidates...)
}

func matchCall(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{Parentheses, Dot, SquareBrackets}

	matched := cursor.MatchAny(candidates...)
	switch matched.Code {
	case dotToken:
		return matchSelector(cursor)
	case parenthesesToken:
		id := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(id[1:len(id)-1]), 0)
		call, err := matchFunctionCall(newCursor)
		if err != nil {
			return nil, err
		}

		call.X, err = matchCall(cursor)
		if err != nil {
			return nil, err
		}
		return call, nil
	}

	return nil, nil
}

func parseIdentity(cursor *parsly.Cursor) (*aexpr.Select, error) {
	candidates := []*parsly.Token{Selector, SelectorBlock}
	matched := cursor.MatchAny(candidates...)
	selectorId := matched.Text(cursor)
	switch matched.Code {
	case parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	case parsly.EOF:
		return &aexpr.Select{ID: selectorId}, nil
	case selectorBlockToken:
		newCursor := parsly.NewCursor("", []byte(selectorId[1:len(selectorId)-1]), 0)
		return parseIdentity(newCursor)
	case selectorToken:
		selector := &aexpr.Select{ID: utils.UpperCaseFirstLetter(selectorId)}
		var err error
		selector.X, err = matchCall(cursor)
		if err != nil {
			return nil, err
		}
		return selector, nil
	}
	return nil, cursor.NewError(candidates...)
}

func matchFunctionCall(cursor *parsly.Cursor) (*aexpr.Call, error) {
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
			expressions = append(expressions, &aexpr.Unary{
				Token: token,
				X:     expression,
			})
		}

	}

	return &aexpr.Call{Args: expressions}, nil
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
