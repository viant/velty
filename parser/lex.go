package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/parsly/matcher"
	vMatcher "github.com/viant/velty/parser/matcher"
)

const (
	specialSignToken = iota
	whiteSpaceToken

	selectorStartToken
	selectorBlockToken
	selectorToken
	complexSelectorToken
	variableNameToken

	ifToken
	elseIfToken
	elseToken
	setToken
	forEachToken
	forToken
	appendToken
	endToken

	inToken

	parentheses
	whitespaceOnlyToken
	expressionBlockToken
	expressionStartToken
	expressionEndToken

	stringToken
	booleanToken
	numberMatcher

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
)

var WhiteSpace = parsly.NewToken(whiteSpaceToken, "Whitespace", matcher.NewWhiteSpace())
var SpecialSign = parsly.NewToken(specialSignToken, "Special sign", vMatcher.NewVelty(true, '#', '$'))

var SelectorBlock = parsly.NewToken(selectorBlockToken, "Selector block", matcher.NewBlock('{', '}', '\\'))
var Selector = parsly.NewToken(selectorToken, "Selector", vMatcher.NewIdentity(false, true))
var ComplexSelector = parsly.NewToken(complexSelectorToken, "Complex selector", vMatcher.NewIdentity(true, false))

var NewVariable = parsly.NewToken(variableNameToken, "New variable", vMatcher.NewIdentity(false, false))
var SelectorStart = parsly.NewToken(selectorStartToken, "Selector start", matcher.NewRunes([]rune{'$'}))

var If = parsly.NewToken(ifToken, "If", matcher.NewFragment("if"))
var ElseIf = parsly.NewToken(elseIfToken, "Else if", matcher.NewFragment("elseif"))
var Else = parsly.NewToken(elseToken, "Else", matcher.NewFragment("else"))
var Set = parsly.NewToken(setToken, "Set", matcher.NewFragment("set"))
var ForEach = parsly.NewToken(forEachToken, "ForEach", matcher.NewFragment("foreach"))
var For = parsly.NewToken(forToken, "For", matcher.NewFragment("for"))
var In = parsly.NewToken(inToken, "In", matcher.NewFragment("in"))
var End = parsly.NewToken(endToken, "End", matcher.NewFragment("end"))

var Parentheses = parsly.NewToken(parentheses, "Parentheses", matcher.NewBlock('(', ')', '\\'))
var WhitespaceOnly = parsly.NewToken(whitespaceOnlyToken, "Whitespace only", vMatcher.NewWhitespaceOnly())
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

var String = parsly.NewToken(stringToken, "String", vMatcher.NewStringMatcher('"'))
var Boolean = parsly.NewToken(booleanToken, "Boolean", matcher.NewFragments([]byte("true"), []byte("false")))
var Number = parsly.NewToken(numberMatcher, "Number", matcher.NewNumber())

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

var Coma = parsly.NewToken(comaToken, "Coma", matcher.NewByte(','))
var NewLine = parsly.NewToken(newLineToken, "New line", vMatcher.NewNewLine())
