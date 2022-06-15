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

type newLine struct {
}

func (t *newLine) Match(cursor *parsly.Cursor) (matched int) {
	for i := cursor.Pos; i < cursor.InputSize; i++ {
		matched++
		if cursor.Input[i] == '\n' {
			return matched
		}
	}
	return matched
}

func NewNewLine() *newLine {
	return &newLine{}
}
