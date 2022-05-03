package parser

import (
	"github.com/viant/parsly"
	ast2 "github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
)

func normalizeTokensIfNeeded(variable ast2.Expression, token ast2.Token, rightOperand ast2.Expression, matched *parsly.TokenMatch) (ast2.Token, ast2.Expression) {
	switch matched.Code {
	case mulEqualToken:
		token = ast2.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast2.MUL, rightOperand)
	case quoEqualToken:
		token = ast2.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast2.QUO, rightOperand)
	case addEqualToken:
		token = ast2.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast2.ADD, rightOperand)
	case subEqualToken:
		token = ast2.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast2.SUB, rightOperand)
	}

	return token, rightOperand
}

func matchToken(matched *parsly.TokenMatch) ast2.Token {
	var token ast2.Token
	switch matched.Code {
	case equalToken:
		token = ast2.EQ
	case greaterToken:
		token = ast2.GTR
	case lessToken:
		token = ast2.LSS
	case lessEqualToken:
		token = ast2.LEQ
	case greaterEqualToken:
		token = ast2.GTE
	case notEqualToken:
		token = ast2.NEQ
	case orToken:
		token = ast2.OR
	case andToken:
		token = ast2.AND
	case addToken:
		token = ast2.ADD
	case subToken:
		token = ast2.SUB
	case mulToken:
		token = ast2.MUL
	case quoToken:
		token = ast2.QUO
	case assignToken:
		token = ast2.ASSIGN
	case negationToken:
		token = ast2.NEG
	}
	return token
}
