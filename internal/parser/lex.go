package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/parsly/matcher"
	matcher2 "github.com/viant/velty/internal/parser/matcher"
)

const (
	specialSignToken = iota
	whiteSpaceToken

	selectorStartToken
	selectorBlockToken
	selectorToken

	ifToken
	elseIfToken
	elseToken
	setToken
	forEachToken
	forToken
	appendToken
	evaluateToken
	endToken

	inToken

	parenthesesToken
	squareBracketsToken
	expressionBlockToken
	expressionStartToken
	expressionEndToken

	stringToken
	booleanToken
	numberToken

	greaterToken
	greaterEqualToken
	lessToken
	lessEqualToken

	equalToken
	notEqualToken
	negationToken

	assignToken

	andToken
	orToken

	addToken
	addEqualToken
	subToken
	subEqualToken
	mulToken
	mulEqualToken
	quoToken
	quoEqualToken

	decrementToken
	incrementToken

	binaryExpressionStartToken

	comaToken
	newLineToken
	dotToken
)

var WhiteSpace = parsly.NewToken(whiteSpaceToken, "Whitespace", matcher.NewWhiteSpace())
var SpecialSign = parsly.NewToken(specialSignToken, "Special sign", matcher2.NewVelty(true, '#', '$'))

var SelectorBlock = parsly.NewToken(selectorBlockToken, "Sel block", matcher.NewBlock('{', '}', '\\'))
var Selector = parsly.NewToken(selectorToken, "Sel", matcher2.NewIdentity())

var SelectorStart = parsly.NewToken(selectorStartToken, "Sel start", matcher.NewRunes([]rune{'$'}))

var If = parsly.NewToken(ifToken, "If", matcher.NewFragment("if"))
var ElseIf = parsly.NewToken(elseIfToken, "Else if", matcher.NewFragment("elseif"))
var Else = parsly.NewToken(elseToken, "Else", matcher.NewFragment("else"))
var Set = parsly.NewToken(setToken, "Set", matcher.NewFragment("set"))
var ForEach = parsly.NewToken(forEachToken, "ForEach", matcher.NewFragment("foreach"))
var For = parsly.NewToken(forToken, "For", matcher.NewFragment("for"))
var In = parsly.NewToken(inToken, "In", matcher.NewFragment("in"))
var Evaluate = parsly.NewToken(evaluateToken, "Evaluate", matcher.NewFragment("evaluate"))
var End = parsly.NewToken(endToken, "End", matcher.NewFragment("end"))

var Parentheses = parsly.NewToken(parenthesesToken, "Parentheses", matcher.NewBlock('(', ')', '\\'))
var SquareBrackets = parsly.NewToken(squareBracketsToken, "Square brackets", matcher.NewBlock('[', ']', '\\'))

var ExpressionBlock = parsly.NewToken(expressionBlockToken, "Expression block", matcher.NewBlock(';', ';', '\\'))
var ExpressionStart = parsly.NewToken(expressionStartToken, "Expression start", matcher.NewByte(';'))
var ExpressionEnd = parsly.NewToken(expressionEndToken, "Expression end", matcher.NewTerminator(';', false))

var Equal = parsly.NewToken(equalToken, "Equal", matcher.NewFragment("=="))
var NotEqual = parsly.NewToken(notEqualToken, "Not equal", matcher.NewFragment("!="))
var Negation = parsly.NewToken(negationToken, "Negation", matcher.NewByte('!'))

var Assign = parsly.NewToken(assignToken, "Assign", matcher.NewFragment("="))

var Greater = parsly.NewToken(greaterToken, "Greater", matcher.NewByte('>'))
var GreaterEqual = parsly.NewToken(greaterEqualToken, "Greater or equal", matcher.NewFragment(">="))
var Less = parsly.NewToken(lessToken, "Less", matcher.NewByte('<'))
var LessEqual = parsly.NewToken(lessEqualToken, "Less or equal", matcher.NewFragment("<="))

var And = parsly.NewToken(andToken, "And", matcher.NewFragment("&&"))
var Or = parsly.NewToken(orToken, "Or", matcher.NewFragment("||"))

var String = parsly.NewToken(stringToken, "String", matcher2.NewStringMatcher('"'))
var Boolean = parsly.NewToken(booleanToken, "Boolean", matcher.NewFragments([]byte("true"), []byte("false")))
var Number = parsly.NewToken(numberToken, "Number", matcher.NewNumber())

var Add = parsly.NewToken(addToken, "Add", matcher.NewByte('+'))
var AddEqual = parsly.NewToken(addEqualToken, "Add equal", matcher.NewBytes([]byte("+=")))
var Sub = parsly.NewToken(subToken, "Subtract", matcher.NewByte('-'))
var SubEqual = parsly.NewToken(subEqualToken, "Subtract equal", matcher.NewBytes([]byte("-=")))
var Multiply = parsly.NewToken(mulToken, "Multiply", matcher.NewByte('*'))
var MultiplyEqual = parsly.NewToken(mulEqualToken, "Multiply equal", matcher.NewBytes([]byte("*=")))
var Quo = parsly.NewToken(quoToken, "Quo", matcher.NewByte('/'))
var QuoEqual = parsly.NewToken(quoEqualToken, "Quo equal", matcher.NewBytes([]byte("/=")))

var Decrement = parsly.NewToken(decrementToken, "Decrement", matcher.NewBytes([]byte("--")))
var Increment = parsly.NewToken(incrementToken, "Increment", matcher.NewBytes([]byte("++")))

var ComaTerminator = parsly.NewToken(comaToken, "Coma", matcher.NewTerminator(',', true))
var NewLine = parsly.NewToken(newLineToken, "New line", matcher2.NewNewLine())
var Dot = parsly.NewToken(dotToken, "Dot", matcher.NewByte('.'))
