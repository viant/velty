package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
)

func matchEquationExpression(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{Parentheses}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	switch matched.Code {
	case negationToken:
		unary := &expr.Unary{Token: ast.NEG}
		expression, err := matchEquationExpression(cursor)
		if err != nil {
			return nil, err
		}
		unary.X = expression
		return unary, nil
	case parentheses:
		expressionValue := matched.Text(cursor)

		matched = cursor.MatchAfterOptional(WhiteSpace, Negation)
		shouldNegate := matched.Code == negationToken

		expressionCursor := parsly.NewCursor("", []byte(expressionValue[1:len(expressionValue)-1]), 0)
		expression, err := matchEquationExpression(expressionCursor)
		if err != nil {
			return nil, err
		}
		err = addEquationIfNeeded(cursor, &expression)
		if err != nil {
			return nil, err
		}

		var result ast.Expression = &expr.Parentheses{P: expression}
		if shouldNegate {
			result = &expr.Unary{
				Token: ast.NEG,
				X:     result,
			}
		}

		return result, nil
	default:
		_, expression, err := matchOperand(cursor, dataTypeMatchers...)
		if err != nil {
			return nil, err
		}

		return expression, nil
	}

}

func matchExpressionBlock(cursor *parsly.Cursor) (*parsly.Cursor, error) {
	expressionMatch := cursor.MatchAfterOptional(WhiteSpace, Parentheses)
	if expressionMatch.Code == parsly.EOF || expressionMatch.Code == parsly.Invalid {
		return nil, cursor.NewError(Parentheses)
	}
	expressionValue := expressionMatch.Text(cursor)
	expressionCursor := parsly.NewCursor("", []byte(expressionValue[1:len(expressionValue)-1]), 0)
	return expressionCursor, nil
}
