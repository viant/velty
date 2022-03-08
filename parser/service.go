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
	match = cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch match.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	case ifToken:
		match = cursor.MatchAfterOptional(WhiteSpace, ExpressionBlock)
		if match.Code == parsly.EOF || match.Code == parsly.Invalid {
			return nil, cursor.NewError(ExpressionBlock)
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

//TODO: Implement #end, #else, handling statements
func matchIf(cursor *parsly.Cursor) (*stmt.If, error) {
	expression, err := matchIfExpression(cursor)
	if err != nil {
		return nil, err
	}

	return &stmt.If{
		Condition: expression,
		Body:      stmt.Block{},
		Else:      nil,
	}, nil

}

func matchIfExpression(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{Negation, ExpressionBlock}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var expression ast.Expression
	var err error
	if isUnaryMatched(matched) {
		expression, err = matchUnaryExpression(cursor, matched)
	} else if matched.Code == expressionBlockToken {
		expressionValue := matched.Text(cursor)
		expressionCursor := parsly.NewCursor("", []byte(expressionValue[1:len(expressionValue)-1]), 0)
		expression, err = matchIfExpression(expressionCursor)
	} else {
		expression, err = matchBinaryExpression(cursor, matched)
	}

	if err != nil {
		return nil, err
	}

	candidates = []*parsly.Token{And, Or}
	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case andToken:
		return matchExpressionCombination(cursor, expression, ast.AND)
	case orToken:
		return matchExpressionCombination(cursor, expression, ast.OR)
	case parsly.EOF:
		return expression, nil
	default:
		return nil, cursor.NewError(candidates...)
	}
}

func matchExpressionCombination(cursor *parsly.Cursor, expression ast.Expression, token ast.Token) (ast.Expression, error) {
	rightExpression, err := matchIfExpression(cursor)
	if err != nil {
		return nil, err
	}

	return &expr.Binary{
		X:     expression,
		Token: token,
		Y:     rightExpression,
	}, nil
}

func matchBinaryExpression(cursor *parsly.Cursor, matched *parsly.TokenMatch) (ast.Expression, error) {
	operandCandidates := []*parsly.Token{String, Number, Boolean}

	leftSideMatcher, leftOperand, err := matchOperand(cursor, operandCandidates...)
	if err != nil {
		return nil, err
	}

	tokenCandidates := []*parsly.Token{NotEqual, Equal, GreaterEqual, Greater, LessEqual, Less, And, Or, ExpressionEnd}
	matched = cursor.MatchAfterOptional(WhiteSpace, tokenCandidates...)
	switch matched.Code {
	case parsly.EOF, expressionEndToken:
		return &expr.Binary{
			X:     leftOperand,
			Token: ast.EQ,
			Y:     expr.BoolExpression("true"),
		}, nil
	}

	token := matchExpressionToken(matched)
	if token == "" {
		return nil, cursor.NewError(tokenCandidates...)
	}

	if leftSideMatcher.Code != selectorToken {
		operandCandidates = []*parsly.Token{leftSideMatcher}
	}

	_, rightOperand, err := matchOperand(cursor, operandCandidates...)
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
		_, expression, err := matchOperand(cursor, Boolean)
		if err != nil {
			return nil, err
		}

		return &expr.Unary{
			Token: ast.NEG,
			X:     expression,
		}, nil
	}

	return nil, cursor.NewError(Boolean, Selector)
}

func matchExpressionToken(matched *parsly.TokenMatch) ast.Token {
	var token ast.Token
	switch matched.Code {
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
	case orToken:
		token = ast.OR
	case andToken:
		token = ast.AND
	}
	return token
}

func matchOperand(cursor *parsly.Cursor, candidates ...*parsly.Token) (*parsly.Token, ast.Expression, error) {
	candidates = append([]*parsly.Token{SelectorStart}, candidates...)

	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var matcher *parsly.Token
	var expression ast.Expression
	var err error

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, nil, cursor.NewError(candidates...)
	case stringToken:
		value := matched.Text(cursor)
		matcher = String
		expression = expr.StringExpression(value[1 : len(value)-1])
	case selectorStartToken:
		matched = cursor.MatchOne(SelectorBlock)
		if matched.Code == parsly.EOF || matched.Code == parsly.Invalid {
			return nil, nil, cursor.NewError(Selector)
		}

		selector := matched.Text(cursor)
		selectorCursor := parsly.NewCursor("", []byte(selector[1:len(selector)-1]), 0)
		expression, err = parseSelector(selectorCursor)
		if err != nil {
			return nil, nil, err
		}
		matcher = Selector

	case numberMatcher:
		value := matched.Text(cursor)
		matcher = Number
		expression = expr.NumberExpression(value)

	case booleanToken:
		value := matched.Text(cursor)
		matcher = Boolean
		expression = expr.BoolExpression(value)
	}

	if matched != nil && matched.Code != booleanToken {
		err = addEquationIfNeeded(cursor, &expression, matcher)
		if err != nil {
			return nil, nil, err
		}

	}
	return matcher, expression, nil
}

func addEquationIfNeeded(cursor *parsly.Cursor, expression *ast.Expression, expressionMatcher *parsly.Token) error {
	candidates := []*parsly.Token{Add, Sub, Multiply, Quo}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	switch matched.Code {
	case parsly.EOF, binaryExpressionStartToken, parsly.Invalid:
		return nil
	}

	token := matchEquationToken(matched)
	switch actual := (*expression).(type) {
	case *expr.Literal:
		_, equationExpression, err := matchOperand(cursor, expressionMatcher)
		if err != nil {
			return err
		}

		*expression = &expr.Binary{
			X:     actual,
			Token: token,
			Y:     equationExpression,
		}
	}

	return nil
}

func matchEquationToken(matched *parsly.TokenMatch) ast.Token {
	var token ast.Token
	switch matched.Code {
	case addToken:
		token = ast.ADD
	case subToken:
		token = ast.SUB
	case mulToken:
		token = ast.MUL
	case quoToken:
		token = ast.QUO
	}

	return token
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
