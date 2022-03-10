package parser

import (
	"fmt"
	"github.com/viant/parsly"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/ast/stmt"
)

func Parse(input []byte) (*stmt.Block, error) {
	if len(input) == 0 {
		return nil, nil
	}

	builder := NewBuilder()
	var tokenMatch *parsly.TokenMatch
	cursor := parsly.NewCursor("", input, 0)

	for cursor.Pos < len(input) {
		tokenMatch = cursor.MatchOne(SpecialSign)
		text := tokenMatch.Text(cursor)

		if tokenMatch.Code == parsly.EOF || cursor.Pos >= len(input) {
			if err := builder.PushStatement(appendToken, stmt.NewAppend(text)); err != nil {
				return nil, err
			}
			break
		}

		if cursor.Input[cursor.Pos] == '#' {
			cursor.MatchOne(NewLine)
			continue
		}

		err := appendStatementIfNeeded(text, builder)
		if err != nil {
			return nil, err
		}

		switch cursor.Input[cursor.Pos-1] {
		case '$':
			statement, err := matchSelector(cursor)
			if err != nil {
				return nil, err
			}
			builder.appendStatement(statement)

		case '#':
			statement, match, err := matchStatement(cursor)
			if err != nil {
				return nil, err
			}

			if err = builder.PushStatement(match, statement); err != nil {
				return nil, err
			}
		}
	}

	if builder.BufferSize() != 0 {
		return nil, fmt.Errorf("unterminated statements on the stack: %v", builder.buffer)
	}

	return builder.Block(), nil
}

func appendStatementIfNeeded(text string, stack *Builder) error {
	text = text[:len(text)-1]
	if len(text) == 0 {
		return nil
	}

	if err := stack.PushStatement(appendToken, stmt.NewAppend(text)); err != nil {
		return err
	}
	return nil
}

func addIfExpression(node ast.Node, expression ast.Node) error {
	switch nodeType := node.(type) {
	case stmt.Condition:
		switch exprType := expression.(type) {
		case *stmt.If:
			nodeType.AddCondition(exprType)
			return nil
		default:
			return fmt.Errorf("expected stmt.If but got %T", expression)
		}
	}
	return fmt.Errorf("expected stmt.Condition but got %T", node)
}

func matchStatement(cursor *parsly.Cursor) (ast.Statement, int, error) {
	candidates := []*parsly.Token{If, ElseIf, Else, Set, ForEach, For, End}
	expressionMatch := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	expressionCode := expressionMatch.Code

	switch expressionMatch.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, 0, cursor.NewError(candidates...)
	case ifToken, elseIfToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		ifStmt, err := matchIf(expressionCursor)
		if err != nil {
			return nil, 0, err
		}
		return ifStmt, expressionCode, nil
	case elseToken:
		return &stmt.If{
			Condition: &expr.Binary{
				X:     expr.BoolExpression("true"),
				Token: "==",
				Y:     expr.BoolExpression("true"),
			},
			Body: stmt.Block{},
		}, expressionCode, nil

	case setToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		assignStmt, err := matchAssign(expressionCursor)
		if err != nil {
			return nil, expressionCode, err
		}

		return assignStmt, expressionCode, nil
	case forEachToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		forEachStmt, err := matchForEach(expressionCursor)
		if err != nil {
			return nil, 0, err
		}

		return forEachStmt, expressionCode, nil

	case forToken:
		expressionCursor, err := matchExpressionBlock(cursor)
		if err != nil {
			return nil, 0, err
		}

		forStmt, err := matchFor(expressionCursor)
		if err != nil {
			return nil, 0, err
		}

		return forStmt, expressionCode, nil

	case endToken:
		return nil, expressionCode, nil
	}

	return nil, 0, cursor.NewError(candidates...)
}

func matchFor(cursor *parsly.Cursor) (*stmt.Range, error) {
	initCursor := extractForSegment(cursor)
	forInit, err := matchAssign(initCursor)
	if err != nil {
		return nil, err
	}

	conditionCursor := extractForSegment(cursor)
	forCondition, err := matchBooleanExpression(conditionCursor)
	if err != nil {
		return nil, err
	}

	forPostCursor := extractForSegment(cursor)
	forPost, err := matchForPost(forPostCursor)
	if err != nil {
		return nil, err
	}

	return &stmt.Range{
		Init: forInit,
		Cond: forCondition,
		Post: forPost,
	}, nil
}

func matchForPost(cursor *parsly.Cursor) (ast.Statement, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}

	candidates := []*parsly.Token{Increment, Decrement, AddEqual, Add, SubEqual, Sub, MultiplyEqual, Multiply, QuoEqual, Quo}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(candidates...)

	case decrementToken:
		return &stmt.Statement{
			X:  variable,
			Op: ast.ASSIGN,
			Y:  expr.BinaryExpression(variable, ast.SUB, expr.NumberExpression("1")),
		}, nil

	case incrementToken:
		return &stmt.Statement{
			X:  variable,
			Op: ast.ASSIGN,
			Y:  expr.BinaryExpression(variable, ast.ADD, expr.NumberExpression("1")),
		}, nil
	}

	token := matchToken(matched)
	_, rightOperand, err := matchOperand(cursor)
	if err != nil {
		return nil, err
	}

	token, rightOperand = normalizeTokensIfNeeded(variable, token, rightOperand, matched)

	return &stmt.Statement{
		X:  variable,
		Op: token,
		Y:  rightOperand,
	}, nil
}

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

func extractForSegment(cursor *parsly.Cursor) *parsly.Cursor {
	candidates := []*parsly.Token{ExpressionBlock, ExpressionStart, ExpressionEnd}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case expressionBlockToken:
		segment := matched.Text(cursor)
		return parsly.NewCursor("", []byte(segment[1:len(segment)-1]), 0)
	case expressionStartToken:
		expressionStartMatch := cursor.MatchAfterOptional(WhiteSpace, ExpressionEnd)
		if expressionStartMatch.Code == parsly.EOF || expressionStartMatch.Code == parsly.Invalid {
			return parsly.NewCursor("", cursor.Input[cursor.Pos:], 0)
		}

		text := expressionStartMatch.Text(cursor)
		return parsly.NewCursor("", []byte(text[:len(text)-1]), 0)
	case expressionEndToken:
		segment := matched.Text(cursor)
		return parsly.NewCursor("", []byte(segment), 0)
	}

	return parsly.NewCursor("", cursor.Input[cursor.Pos:], 0)
}

func matchForEach(cursor *parsly.Cursor) (*stmt.ForEach, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}
	candidates := []*parsly.Token{Coma}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var index *expr.Select
	if matched.Code == comaToken {
		index, err = matchVariable(cursor)
		if err != nil {
			return nil, err
		}
	}
	candidates = []*parsly.Token{In}

	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)
	if matched.Code == parsly.Invalid || matched.Code == parsly.EOF {
		return nil, cursor.NewError(candidates...)
	}

	dataSet, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}

	return &stmt.ForEach{
		Index: index,
		Item:  variable,
		Set:   dataSet,
		Body:  stmt.Block{},
	}, nil
}

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

func matchExpressionBlock(cursor *parsly.Cursor) (*parsly.Cursor, error) {
	expressionMatch := cursor.MatchAfterOptional(WhiteSpace, Parentheses)
	if expressionMatch.Code == parsly.EOF || expressionMatch.Code == parsly.Invalid {
		return nil, cursor.NewError(Parentheses)
	}
	expressionValue := expressionMatch.Text(cursor)
	expressionCursor := parsly.NewCursor("", []byte(expressionValue[1:len(expressionValue)-1]), 0)
	return expressionCursor, nil
}

func matchAssign(cursor *parsly.Cursor) (*stmt.Statement, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}

	tokenCandidates := []*parsly.Token{Assign}
	matched := cursor.MatchAfterOptional(WhiteSpace, tokenCandidates...)
	if matched.Code == parsly.EOF || matched.Code == parsly.Invalid {
		return nil, cursor.NewError(tokenCandidates...)
	}

	token := matchToken(matched)
	if token == "" {
		return nil, fmt.Errorf("didn't found operator token for given token %v", matched.Name)
	}

	_, expression, err := matchOperand(cursor, Boolean, String, Number)
	if err != nil {
		return nil, err
	}

	return &stmt.Statement{
		X:  variable,
		Op: token,
		Y:  expression,
	}, nil
}

func matchIf(cursor *parsly.Cursor) (*stmt.If, error) {
	expression, err := matchBooleanExpression(cursor)
	if err != nil {
		return nil, err
	}

	return &stmt.If{
		Condition: expression,
		Body:      stmt.Block{},
		Else:      nil,
	}, nil
}

func matchBooleanExpression(cursor *parsly.Cursor) (ast.Expression, error) {
	candidates := []*parsly.Token{Negation, Parentheses}
	matched := cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var expression ast.Expression
	var err error
	if isUnaryMatched(matched) {
		expression, err = matchUnaryExpression(cursor, matched)
	} else if matched.Code == parentheses {
		expressionValue := matched.Text(cursor)
		expressionCursor := parsly.NewCursor("", []byte(expressionValue[1:len(expressionValue)-1]), 0)
		expression, err = matchBooleanExpression(expressionCursor)
	} else {
		expression, err = matchBinaryExpression(cursor, matched)
	}

	if err != nil {
		return nil, err
	}

	candidates = []*parsly.Token{And, Or, Equal, NotEqual}
	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)
	switch matched.Code {
	case andToken:
		return matchExpressionCombination(cursor, expression, ast.AND)
	case orToken:
		return matchExpressionCombination(cursor, expression, ast.OR)
	case parsly.EOF:
		return expression, nil
	case parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	default:
		token := matchToken(matched)
		return matchExpressionCombination(cursor, expression, token)
	}
}

func matchExpressionCombination(cursor *parsly.Cursor, expression ast.Expression, token ast.Token) (ast.Expression, error) {
	rightExpression, err := matchBooleanExpression(cursor)
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

	tokenCandidates := []*parsly.Token{NotEqual, Equal, GreaterEqual, Greater, LessEqual, Less, And, Or, WhitespaceOnly}
	matched = cursor.MatchAfterOptional(WhiteSpace, tokenCandidates...)
	switch matched.Code {
	case parsly.EOF, whitespaceOnlyToken:
		return &expr.Binary{
			X:     leftOperand,
			Token: ast.EQ,
			Y:     expr.BoolExpression("true"),
		}, nil
	}

	token := matchToken(matched)
	if token == "" {
		return nil, cursor.NewError(tokenCandidates...)
	}

	if leftSideMatcher.Code != complexSelectorToken && leftSideMatcher.Code != selectorToken {
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

	return nil, cursor.NewError(Boolean, ComplexSelector)
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
		expression, err = matchSelector(cursor)
		if err != nil {
			return nil, nil, err
		}

		matcher = ComplexSelector

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

	token := matchToken(matched)
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
		selector, err := parseIdentity(selectorCursor, true)
		if err != nil {
			return nil, err
		}
		return selector, nil
	}

	return nil, cursor.NewError(candidates...)
}

func parseIdentity(cursor *parsly.Cursor, fullMatch bool) (*expr.Select, error) {
	var candidates []*parsly.Token
	if fullMatch {
		candidates = []*parsly.Token{ComplexSelector}
	} else {
		candidates = []*parsly.Token{NewVariable}
	}

	matched := cursor.MatchAny(candidates...)
	id := matched.Text(cursor)

	switch matched.Code {
	case parsly.EOF, parsly.Invalid:
		return nil, cursor.NewError(candidates...)
	}

	candidates = []*parsly.Token{Parentheses}
	matched = cursor.MatchAfterOptional(WhiteSpace, candidates...)

	var call *expr.Call
	if matched.Code == parentheses {
		callValue := matched.Text(cursor)
		callCursor := parsly.NewCursor("", []byte(callValue[1:len(callValue)-1]), 0)

		var err error
		call, err = matchFunctionCall(callCursor)
		if err != nil {
			return nil, err
		}
	}

	return &expr.Select{
		ID:   id,
		Call: call,
	}, nil
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
