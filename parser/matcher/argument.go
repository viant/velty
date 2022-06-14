package matcher

import (
	"github.com/viant/parsly"
)

type argument struct {
}

func (a *argument) Match(cursor *parsly.Cursor) (matched int) {
	depth := 0
	for i := cursor.Pos; i < cursor.InputSize; i++ {
		matched++
		switch cursor.Input[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				return matched
			}
		}
	}
	return matched
}

func NewArgumentMatcher() *argument {
	return &argument{}
}
