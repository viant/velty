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

//TODO: Refactor
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
		selector, err := parseIdentity(selectorCursor, false)
		if err != nil {
			return nil, err
		}

		var call *expr.Call
		blockCursor, err := matchExpressionBlock(cursor)
		if err == nil {
			call, _ = matchFunctionCall(blockCursor)
			if call != nil {
				selector.X = call
			}
		}

		matched = cursor.MatchOne(Dot)

		if matched.Code == dotToken {
			if call != nil {
				call.X, err = matchSelector(cursor)
				if err != nil {
					return nil, err
				}
			} else {
				selector.X, err = matchSelector(cursor)
				if err != nil {
					return nil, err
				}
			}
		}
		return selector, nil
	}

	return nil, cursor.NewError(candidates...)
}

//TODO: Refactor
func parseIdentity(cursor *parsly.Cursor, fullMatch bool) (*expr.Select, error) {
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
		return parseIdentity(newCursor, true)
	}

	candidates = []*parsly.Token{Dot, Parentheses, SquareBrackets}
	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)
	selector := &expr.Select{ID: selectorId}
	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return selector, nil
	case dotToken:
		call, err := parseIdentity(cursor, fullMatch)
		if err != nil {
			return nil, err
		}
		selector.X = call

	case parenthesesToken:
		text := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(text[1:len(text)-1]), 0)
		call, err := matchFunctionCall(newCursor)
		if err != nil {
			return nil, err
		}
		selector.X = call
		if matchNextIdentity(cursor) {
			call.X, err = parseIdentity(cursor, fullMatch)
			if err != nil {
				return nil, err
			}
		}
	case squareBracketsToken:
		text := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(text[1:len(text)-1]), 0)
		_, operandExpr, err := matchOperand(newCursor, dataTypeMatchers...)
		if err != nil {
			return nil, err
		}

		call := &expr.SliceIndex{X: operandExpr}
		if matchNextIdentity(cursor) {
			call.Y, err = parseIdentity(cursor, fullMatch)
			if err != nil {
				return nil, err
			}
		}

		selector.X = call
		return selector, nil
	}

	if fullMatch {
		cursor.MatchOne(WhiteSpace)
		if cursor.Pos != cursor.InputSize {
			return nil, cursor.NewError(WhiteSpace)
		}
	}
	return selector, nil
}

func matchNextIdentity(cursor *parsly.Cursor) bool {
	return cursor.MatchAfterOptional(WhiteSpace, Dot).Code == dotToken
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
