package matcher

import (
	"github.com/viant/parsly"
)

type whitespaceOnly struct {
}

func (t *whitespaceOnly) Match(cursor *parsly.Cursor) (matched int) {
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

func NewWhitespaceOnly() *whitespaceOnly {
	return &whitespaceOnly{}
}
