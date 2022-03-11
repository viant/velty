package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
)

func normalizeTokensIfNeeded(variable ast.Expression, token ast.Token, rightOperand ast.Expression, matched *parsly.TokenMatch) (ast.Token, ast.Expression) {
	switch matched.Code {
	case mulEqualToken:
		token = ast.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast.MUL, rightOperand)
	case quoEqualToken:
		token = ast.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast.QUO, rightOperand)
	case addEqualToken:
		token = ast.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast.ADD, rightOperand)
	case subEqualToken:
		token = ast.ASSIGN
		rightOperand = expr.BinaryExpression(variable, ast.SUB, rightOperand)
	}

	return token, rightOperand
}

func matchToken(matched *parsly.TokenMatch) ast.Token {
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
	case addToken:
		token = ast.ADD
	case subToken:
		token = ast.SUB
	case mulToken:
		token = ast.MUL
	case quoToken:
		token = ast.QUO
	case assignToken:
		token = ast.ASSIGN
	case negationToken:
		token = ast.NEG
	}
	return token
}
