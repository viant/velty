package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
)

func Parse(input []byte) (ast.Node, error) {
	var tokenMatch *parsly.TokenMatch
	cursor := parsly.NewCursor("", input, 0)

	tokenMatch = cursor.MatchOne(SpecialSign)
	switch tokenMatch.Code {
	case parsly.EOF:
		return nil, nil
	case specialSignToken:
		switch cursor.Input[cursor.Pos-1] {
		case '$':
			return matchSelector(tokenMatch, cursor)
		case '#':
			return matchExpression(tokenMatch, cursor)
		}
	}
	return nil, cursor.NewError(SpecialSign)
}

func matchExpression(match *parsly.TokenMatch, cursor *parsly.Cursor) (ast.Node, error) {
	candidates := []*parsly.Token{If}
	match = cursor.MatchAny(candidates...)
	switch match.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	case ifToken:
		match = cursor.MatchOne(IfBlock)
		if match.Code == parsly.EOF || match.Code == parsly.Invalid {
			return nil, cursor.NewError(IfBlock)
		}
		ifCondition := match.Text(cursor)
		conditionCursor := parsly.NewCursor("", []byte(ifCondition[1:len(ifCondition)-1]), 0)
		ifStmt, err := matchIf(conditionCursor)
		if err != nil {
			return nil, err
		}
		return ifStmt, nil
	}

	return nil, cursor.NewError(candidates...)
}

//TODO: Implement #end, #else, unary if, condition composition, type checking, handling statements
func matchIf(cursor *parsly.Cursor) (*stmt.If, error) {
	operandCandidates := []*parsly.Token{Negation}
	matched := cursor.MatchAny(operandCandidates...)

	var err error
	var expression ast.Expression
	if matched.Code != parsly.EOF && matched.Code != parsly.Invalid {
		expression, err = matchUnaryExpression(cursor, matched)
	} else {
		expression, err = matchBinaryExpression(cursor, matched)
	}

	if err != nil {
		return nil, err
	}

	return &stmt.If{
		Condition: expression,
		Body:      stmt.Block{},
		Else:      nil,
	}, nil

}

func matchBinaryExpression(cursor *parsly.Cursor, matched *parsly.TokenMatch) (ast.Expression, error) {
	operandCandidates := []*parsly.Token{StringMatcher, NumberMatcher, BooleanMatcher}

	leftSideMatcher, leftOperand, err := matchOperand(cursor, operandCandidates)
	if err != nil {
		return nil, err
	}

	tokenCandidates := []*parsly.Token{NotEqual, Equal, GreaterEqual, Greater, LessEqual, Less}
	matched = cursor.MatchAfterOptional(WhiteSpace, tokenCandidates...)
	token, err := matchToken(cursor, matched, tokenCandidates)
	if err != nil {
		return nil, err
	}

	if leftSideMatcher.Code != selectorToken {
		operandCandidates = []*parsly.Token{Selector, leftSideMatcher}
	}

	_, rightOperand, err := matchOperand(cursor, operandCandidates)
	if err != nil {
		return nil, err
	}
	return &expr.Binary{
		X:     leftOperand,
		Token: token,
		Y:     rightOperand,
	}, nil
}

func matchUnaryExpression(cursor *parsly.Cursor, matched *parsly.TokenMatch) (ast.Expression, error) {
	switch matched.Code {
	case negationToken:
		_, expression, err := matchOperand(cursor, []*parsly.Token{BooleanMatcher})
		if err != nil {
			return nil, err
		}

		return &expr.Unary{
			Token: ast.NEG,
			X:     expression,
		}, nil
	}

	return nil, cursor.NewError(BooleanMatcher, Selector)
}

func matchToken(cursor *parsly.Cursor, matched *parsly.TokenMatch, tokenCandidates []*parsly.Token) (ast.Token, error) {
	var token ast.Token
	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return "", cursor.NewError(tokenCandidates...)
	case equalToken:
		token = ast.EQ
	case greaterToken:
		token = ast.GTR
	case lessToken:
		token = ast.LSS
	case lessEqualToken:
		token = ast.LEQ
	case greaterEqualToken:
		token = ast.GTE
	case notEqualToken:
		token = ast.NEQ
	}
	return token, nil
}

func matchOperand(cursor *parsly.Cursor, candidates []*parsly.Token) (*parsly.Token, ast.Expression, error) {
	candidates = append([]*parsly.Token{SelectorStart}, candidates...)

	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, nil, cursor.NewError(candidates...)
	case stringToken:
		value := matched.Text(cursor)
		return StringMatcher, expr.StringExpression(value[1 : len(value)-1]), nil
	case selectorStartToken:
		matched = cursor.MatchOne(SelectorBlock)
		if matched.Code == parsly.EOF || matched.Code == parsly.Invalid {
			return nil, nil, cursor.NewError(Selector)
		}

		selector := matched.Text(cursor)
		selectorCursor := parsly.NewCursor("", []byte(selector[1:len(selector)-1]), 0)
		operand, err := parseSelector(selectorCursor)
		if err != nil {
			return nil, nil, err
		}
		return Selector, operand, nil
	case numberMatcher:
		value := matched.Text(cursor)
		return NumberMatcher, expr.NumberExpression(value), nil

	case booleanToken:
		value := matched.Text(cursor)
		return BooleanMatcher, expr.BoolExpression(value), nil
	}
	return nil, nil, cursor.NewError(candidates...)
}

func matchSelector(tokenMatch *parsly.TokenMatch, cursor *parsly.Cursor) (ast.Node, error) {
	tokenMatch = cursor.MatchOne(SelectorBlock)
	if tokenMatch.Code == parsly.EOF {
		return nil, cursor.NewError(SelectorBlock)
	}

	ID := tokenMatch.Text(cursor)

	selectorCursor := parsly.NewCursor("", []byte(ID[1:len(ID)-1]), 0)
	selector, err := parseSelector(selectorCursor)
	if err != nil {
		return nil, err
	}
	return selector, nil
}

func parseSelector(cursor *parsly.Cursor) (*expr.Select, error) {
	matched := cursor.MatchOne(Selector)

	id := matched.Text(cursor)
	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(Selector)
	}

	return &expr.Select{ID: id}, nil
}
