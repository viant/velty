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

	ifToken
	ifBlockToken

	stringToken
	booleanToken
	numberMatcher
	floatToken

	greaterToken
	greaterEqualToken
	lessToken
	lessEqualToken

	equalToken
	notEqualToken
	negationToken
)

var WhiteSpace = parsly.NewToken(whiteSpaceToken, "Whitespace", matcher.NewWhiteSpace())
var SpecialSign = parsly.NewToken(specialSignToken, "Special sign", matcher.NewRunes([]rune{'#', '$'}))

var SelectorBlock = parsly.NewToken(selectorBlockToken, "Selector block", matcher.NewBlock('{', '}', '\\'))
var Selector = parsly.NewToken(selectorToken, "Selector", vMatcher.NewIdentity())
var SelectorStart = parsly.NewToken(selectorStartToken, "Selector start", matcher.NewRunes([]rune{'$'}))

var If = parsly.NewToken(ifToken, "If", matcher.NewFragment("if"))
var IfBlock = parsly.NewToken(ifBlockToken, "If block", matcher.NewBlock('(', ')', '\\'))

var Equal = parsly.NewToken(equalToken, "Equal", matcher.NewFragment("=="))
var NotEqual = parsly.NewToken(equalToken, "Not equal", matcher.NewFragment("!="))
var Negation = parsly.NewToken(negationToken, "Negation", matcher.NewByte('!'))

var Greater = parsly.NewToken(greaterToken, "Greater", matcher.NewByte('>'))
var GreaterEqual = parsly.NewToken(greaterEqualToken, "Greater or equal", matcher.NewFragment(">="))
var Less = parsly.NewToken(lessToken, "Less", matcher.NewByte('<'))
var LessEqual = parsly.NewToken(lessEqualToken, "Less or equal", matcher.NewFragment("<="))

var StringMatcher = parsly.NewToken(stringToken, "String", vMatcher.NewStringMatcher('"'))
var BooleanMatcher = parsly.NewToken(booleanToken, "Boolean", matcher.NewFragments([]byte("true"), []byte("false")))
var NumberMatcher = parsly.NewToken(numberMatcher, "Number matcher", matcher.NewNumber())
