package parser

import (
	"fmt"
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/utils"
	"strconv"
	"strings"
)

func matchVariable(cursor *parsly.Cursor) (*expr.Select, error) {
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
		return MatchSelector(cursor)
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

	return &expr.Range{
		X: expr.Number(begin),
		Y: expr.Number(finish),
	}, nil
}

func MatchSelector(cursor *parsly.Cursor) (*expr.Select, error) {
	selectorStart := cursor.Pos
	matched := cursor.MatchOne(Negation) // Java velocity supports `$!`. If String is null, then it will be replaced with Empty String. In Go - we don't need that
	matched = cursor.MatchOne(SelectorBlock)

	if matched.Code == selectorBlockToken {
		ID := matched.Text(cursor)
		selectorCursor := parsly.NewCursor("", []byte(ID[1:len(ID)-1]), 0)
		result, err := MatchSelector(selectorCursor)
		if err != nil {
			return nil, err
		}

		if selectorCursor.Pos < selectorCursor.InputSize {
			return nil, fmt.Errorf("expected to match all data, but couldn't match %v", string(cursor.Input[cursor.Pos:]))
		}

		result.FullName = "$" + ID
		return result, err

	}

	candidates := []*parsly.Token{Selector}
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

		selector.FullName = "$" + string(cursor.Input[selectorStart:cursor.Pos])
		return selector, nil
	}

	return nil, cursor.NewError(candidates...)
}

func matchCall(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{Parentheses, Dot, SquareBrackets}

	matched := cursor.MatchAny(candidates...)
	switch matched.Code {
	case dotToken:
		return MatchSelector(cursor)
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

	case squareBracketsToken:
		id := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(id[1:len(id)-1]), 0)
		_, expression, err := matchOperand(newCursor, Number)
		if err != nil {
			return nil, err
		}

		index := &expr.SliceIndex{
			X: expression,
		}
		index.Y, err = matchCall(cursor)
		if err != nil {
			return nil, err
		}

		return index, nil
	}

	return nil, nil
}

func parseIdentity(cursor *parsly.Cursor) (*expr.Select, error) {
	candidates := []*parsly.Token{Selector, SelectorBlock}
	matched := cursor.MatchAny(candidates...)
	selectorId := matched.Text(cursor)
	switch matched.Code {
	case parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	case parsly.EOF:
		return &expr.Select{ID: selectorId}, nil
	case selectorBlockToken:
		newCursor := parsly.NewCursor("", []byte(selectorId[1:len(selectorId)-1]), 0)
		return parseIdentity(newCursor)
	case selectorToken:
		selector := &expr.Select{ID: utils.UpperCaseFirstLetter(selectorId)}
		var err error
		selector.X, err = matchCall(cursor)
		if err != nil {
			return nil, err
		}
		return selector, nil
	}
	return nil, cursor.NewError(candidates...)
}

func matchFunctionCall(cursor *parsly.Cursor) (*expr.Call, error) {
	expressions := make([]ast.Expression, 0)

	for cursor.Pos < cursor.InputSize-1 {
		argumentCursor := extractArgument(cursor)
		_, expression, err := matchOperand(argumentCursor, String, Boolean, Number)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expression)
	}

	return &expr.Call{Args: expressions}, nil
}

func extractArgument(cursor *parsly.Cursor) *parsly.Cursor {
	candidates := []*parsly.Token{Argument}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case argumentToken:
		text := matched.Text(cursor)
		return parsly.NewCursor("", []byte(text), 0)
	}

	return cursor
}
