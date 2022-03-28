package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/internal/ast"
	aexpr "github.com/viant/velty/internal/ast/expr"
)

var dataTypeMatchers = []*parsly.Token{String, Boolean, Number}

func matchOperand(cursor *parsly.Cursor, candidates ...*parsly.Token) (*parsly.Token, ast.Expression, error) {
	matched := cursor.MatchAfterOptional(WhiteSpace, Negation)
	hasNegation := matched.Code == negationToken

	candidates = append([]*parsly.Token{SelectorStart}, candidates...)

	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var matcher *parsly.Token
	var expression ast.Expression
	var err error

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, nil, cursor.NewError(candidates...)
	case stringToken:
		value := matched.Text(cursor)
		matcher = String
		expression = aexpr.StringExpression(value[1 : len(value)-1])

	case selectorStartToken:
		expression, err = matchSelector(cursor)
		if err != nil {
			return nil, nil, err
		}

		matcher = Selector

	case numberToken:
		value := matched.Text(cursor)
		matcher = Number
		expression = aexpr.NumberExpression(value)

	case booleanToken:
		value := matched.Text(cursor)
		matcher = Boolean
		expression = aexpr.BoolExpression(value)
	}

	if hasNegation {
		expression = &aexpr.Unary{
			Token: ast.NEG,
			X:     expression,
		}
	}
	err = addEquationIfNeeded(cursor, &expression)
	if err != nil {
		return nil, nil, err
	}

	return matcher, expression, nil
}

func addEquationIfNeeded(cursor *parsly.Cursor, expression *ast.Expression) error {
	for {
		candidates := []*parsly.Token{Add, Sub, Multiply, Quo, NotEqual, Negation, Equal, And, Or, GreaterEqual, Greater, LessEqual, Less}
		matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

		switch matched.Code {
		case parsly.EOF, binaryExpressionStartToken, parsly.Invalid:
			return nil
		}

		token := matchToken(matched)

		eprCursor, err := matchExpressionBlock(cursor)

		var rightExpression ast.Expression
		if err == nil {
			rightExpression, err = matchEquationExpression(eprCursor)
			rightExpression = &aexpr.Parentheses{P: rightExpression}
		} else {
			_, rightExpression, err = matchOperand(cursor, dataTypeMatchers...)
		}

		if err != nil {
			return err
		}
		hasPrecedence := isPrecedenceToken(token)

		if hasPrecedence {
			y, ok := rightExpression.(*aexpr.Binary)
			if ok && !isPrecedenceToken(y.Token) {
				expression = adjustPrecedence(expression, token, y)
				continue
			}

		}

		*expression = &aexpr.Binary{
			X:     *expression,
			Token: token,
			Y:     rightExpression,
		}
	}
}

func adjustPrecedence(expression *ast.Expression, token ast.Token, y *aexpr.Binary) *ast.Expression {
	p := &aexpr.Parentheses{}
	p.P = &aexpr.Binary{
		X:     *expression,
		Token: token,
		Y:     y.X,
	}

	*expression = &aexpr.Binary{
		X:     p,
		Token: y.Token,
		Y:     y.Y,
	}
	return expression
}

func isPrecedenceToken(token ast.Token) bool {
	hasPrecedence := token == ast.MUL || token == ast.QUO
	return hasPrecedence
}
