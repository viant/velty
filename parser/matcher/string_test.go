package matcher

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/parsly"
	"testing"
)

func TestStringMatcher_Match(t *testing.T) {
	testcases := []struct {
		input          string
		matched        int
		cursorPosition int
		quoteChar      byte
		description    string
	}{
		{
			description:    "regular string",
			input:          "'abc'",
			matched:        5,
			cursorPosition: 0,
			quoteChar:      '\'',
		},
		{
			description:    "string escaped",
			input:          `'abc\''`,
			matched:        7,
			cursorPosition: 0,
			quoteChar:      '\'',
		},
	}

	for _, testcase := range testcases {
		matcher := NewStringMatcher(testcase.quoteChar)
		cursor := parsly.NewCursor("", []byte(testcase.input), 0)
		cursor.Pos = testcase.cursorPosition
		assert.Equal(t, testcase.matched, matcher.Match(cursor), testcase.description)
	}
}
