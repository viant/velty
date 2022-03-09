package matcher

import "github.com/viant/parsly"

type expressionStart struct {
	values    []byte
	inclusive bool
}

func (t *expressionStart) Match(cursor *parsly.Cursor) (matched int) {
	hasMatch := false
outer:
	for _, c := range cursor.Input[cursor.Pos:] {
		matched++
		for _, terminator := range t.values {
			if hasMatch = c == terminator; hasMatch {
				if !t.inclusive {
					matched--
				}
				break outer
			}
		}

	}
	if matched < 0 {
		matched = 0
	}

	return matched
}

//NewExpression creates a terminator byte matcher
func NewExpression(inclusive bool, values ...byte) *expressionStart {
	return &expressionStart{values: values, inclusive: inclusive}
}
