package ast

type Token string

const (
	ADD        = Token("+")
	SUB        = Token('-')
	MUL        = Token('*')
	QUO        = Token('/')
	GTR        = Token('>')
	GTE        = Token(">=")
	LSS        = Token("<")
	LEQ        = Token("<=")
	EQ         = Token("==")
	NEQ        = Token("!=")
	NEG        = Token("!")
	ASSIGN     = Token("=")
	ADD_ASSIGN = Token("+=")
	INC        = Token("++")
	DEC        = Token("--")
)
