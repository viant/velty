package parser

import (
	"github.com/viant/parsly"
	ast2 "github.com/viant/velty/ast"
	aexpr "github.com/viant/velty/ast/expr"
)

var dataTypeMatchers = []*parsly.Token{String, Boolean, Number}

func matchOperand(cursor *parsly.Cursor, candidates ...*parsly.Token) (*parsly.Token, ast2.Expression, error) {
	matched := cursor.MatchAfterOptional(WhiteSpace, Negation)
	hasNegation := matched.Code == negationToken

	candidates = append([]*parsly.Token{Quote, SelectorStart, Parentheses}, candidates...)

	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var matcher *parsly.Token
	var expression ast2.Expression
	var err error

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, nil, cursor.NewError(candidates...)
	case parenthesesToken:
		text := matched.Text(cursor)
		newCursor := parsly.NewCursor("", []byte(text[1:len(text)-1]), 0)
		token, expr, err := matchOperand(newCursor, candidates...)
		if err != nil {
			return nil, nil, err
		}

		if hasNegation {
			expr = &aexpr.Unary{
				Token: ast2.NEG,
				X:     expr,
			}
		}

		return token, expr, nil
	case stringToken:
		value := matched.Text(cursor)
		matcher = String
		expression = aexpr.StringLiteral(value[1 : len(value)-1])

	case selectorStartToken:
		expression, err = matchSelector(cursor)
		if err != nil {
			return nil, nil, err
		}

		matcher = Selector

	case numberToken:
		value := matched.Text(cursor)
		matcher = Number
		expression = aexpr.NumberLiteral(value)

	case booleanToken:
		value := matched.Text(cursor)
		matcher = Boolean
		expression = aexpr.BoolLiteral(value)

	case quoteToken:
		matched = cursor.MatchOne(StringFinish)
		if matched.Code != stringFinishToken {
			return nil, nil, cursor.NewError(StringFinish)
		}

		value := matched.Text(cursor)
		if len(value) == 1 { // matched `"`
			matcher = String
			expression = aexpr.StringLiteral("")
		} else {
			newCursor := parsly.NewCursor("", []byte(value[:len(value)-1]), 0)

			matcher, expression, err = matchOperand(newCursor, candidates...)
			if err != nil {
				expression = aexpr.StringLiteral(value[:len(value)-1])
			}
		}

	}

	if hasNegation {
		expression = &aexpr.Unary{
			Token: ast2.NEG,
			X:     expression,
		}
	}
	err = addEquationIfNeeded(cursor, &expression)
	if err != nil {
		return nil, nil, err
	}

	return matcher, expression, nil
}

func addEquationIfNeeded(cursor *parsly.Cursor, expression *ast2.Expression) error {
	for {
		candidates := []*parsly.Token{Add, Sub, Multiply, Quo, NotEqual, Negation, Equal, And, Or, GreaterEqual, Greater, LessEqual, Less}
		matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

		switch matched.Code {
		case parsly.EOF, binaryExpressionStartToken, parsly.Invalid:
			return nil
		}

		token := matchToken(matched)

		eprCursor, err := matchExpressionBlock(cursor)

		var rightExpression ast2.Expression
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

func adjustPrecedence(expression *ast2.Expression, token ast2.Token, y *aexpr.Binary) *ast2.Expression {
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

func isPrecedenceToken(token ast2.Token) bool {
	hasPrecedence := token == ast2.MUL || token == ast2.QUO
	return hasPrecedence
}
