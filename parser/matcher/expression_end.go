package matcher

import (
	"github.com/viant/parsly"
)

type expressionEnd struct {
}

func (t *expressionEnd) Match(cursor *parsly.Cursor) (matched int) {
	for i := cursor.Pos; i < cursor.InputSize; i++ {
		switch cursor.Input[i] {
		case ' ', '\n', '\t', '\r', '\v', '\f', 0x85, 0x40:
			matched++
		default:
			return 0
		}
	}
	return matched
}

func NewExpressionEnd() *expressionEnd {
	return &expressionEnd{}
}
