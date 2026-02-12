package parser

import (
	"fmt"
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

		expr = &aexpr.Parentheses{P: expr}
		// span covers parentheses expression including brackets
		pStart := cursor.Pos - len(text)
		pEnd := cursor.Pos - 1
		recordSpan(expr, pStart, pEnd)

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
		s := cursor.Pos - len(value)
		e := cursor.Pos - 1
		recordSpan(expression, s, e)

	case selectorStartToken:
		expression, err = MatchSelector(cursor)
		if err != nil {
			return nil, nil, err
		}

		matcher = Selector

	case numberToken:
		value := matched.Text(cursor)
		matcher = Number
		expression = aexpr.NumberLiteral(value)
		s := cursor.Pos - len(value)
		e := cursor.Pos - 1
		recordSpan(expression, s, e)

	case booleanToken:
		value := matched.Text(cursor)
		matcher = Boolean
		expression = aexpr.BoolLiteral(value)
		s := cursor.Pos - len(value)
		e := cursor.Pos - 1
		recordSpan(expression, s, e)

	case quoteToken:
		matched = cursor.MatchOne(StringFinish)
		if matched.Code != stringFinishToken {
			return nil, nil, cursor.NewError(StringFinish)
		}

		value := matched.Text(cursor)
		if len(value) == 1 { // matched `"`
			matcher = String
			expression = aexpr.StringLiteral("")
			s := cursor.Pos - len(value)
			e := cursor.Pos - 1
			recordSpan(expression, s, e)
		} else {
			newCursor := parsly.NewCursor("", []byte(value[:len(value)-1]), 0)

			matcher, expression, err = matchOperand(newCursor, candidates...)
			if err != nil {
				expression = aexpr.StringLiteral(value[:len(value)-1])
				s := cursor.Pos - len(value)
				e := cursor.Pos - 1
				recordSpan(expression, s, e)
			} else {
				if _, ok := expression.(*aexpr.Select); !ok {
					expression = aexpr.StringLiteral(value[:len(value)-1])
					s := cursor.Pos - len(value)
					e := cursor.Pos - 1
					recordSpan(expression, s, e)
				}
			}
		}

	}

	if hasNegation {
		expression = &aexpr.Unary{
			Token: ast2.NEG,
			X:     expression,
		}
		// approximate span: extend one char to the left of operand
		if s, ok := getSpan(expression.(*aexpr.Unary).X); ok {
			recordSpan(expression, s.Start-1, s.End)
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
		candidates := []*parsly.Token{Add, Sub, Multiply, Quo, NotEqual, Negation, Equal, And, Or, GreaterEqual, Greater, LessEqual, Less, Assign}
		matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

		switch matched.Code {
		case parsly.EOF, binaryExpressionStartToken, parsly.Invalid:
			return nil
		}

		token := matchToken(matched)
		if token == ast2.ASSIGN {
			return fmt.Errorf("assignment in expression is not allowed")
		}

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
		// derive span from children if available
		if lx, okx := getSpan((*expression).(*aexpr.Binary).X); okx {
			if ry, oky := getSpan((*expression).(*aexpr.Binary).Y); oky {
				start := lx.Start
				if ry.Start < start {
					start = ry.Start
				}
				end := lx.End
				if ry.End > end {
					end = ry.End
				}
				recordSpan(*expression, start, end)
			}
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
