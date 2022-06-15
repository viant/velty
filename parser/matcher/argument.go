package matcher

import (
	"github.com/viant/parsly"
)

type argument struct {
}

func (a *argument) Match(cursor *parsly.Cursor) (matched int) {
	depth := 0
	inQuote := false
	for i := cursor.Pos; i < cursor.InputSize; i++ {
		matched++
		switch cursor.Input[i] {
		case '(':
			if !inQuote {
				depth++
			}
		case ')':
			if depth > 0 && !inQuote {
				depth--
			}
		case '"':
			inQuote = !inQuote
		case ',':
			if depth == 0 && !inQuote {
				return matched
			}
		}
	}
	return matched
}

func NewArgumentMatcher() *argument {
	return &argument{}
}
