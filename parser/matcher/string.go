package matcher

import (
	"github.com/viant/parsly"
)

type stringMatcher struct {
	quoteChar byte
}

func (s stringMatcher) Match(cursor *parsly.Cursor) (matched int) {
	input := cursor.Input
	pos := cursor.Pos

	if input[pos] != s.quoteChar {
		return 0
	}
	matched++
	pos++

	escaped := false
	for ; pos < len(input); pos++ {
		matched++
		switch input[pos] {
		case s.quoteChar:
			if !escaped {
				return matched
			}
			escaped = false
		case '\\':
			escaped = !escaped
		default:
			escaped = false
		}
	}

	return 0
}

func NewStringMatcher(quoteChar byte) *stringMatcher {
	return &stringMatcher{quoteChar: quoteChar}
}
