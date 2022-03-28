package matcher

import "github.com/viant/parsly"

type veltyStart struct {
	values    []byte
	inclusive bool
}

func (t *veltyStart) Match(cursor *parsly.Cursor) (matched int) {
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

//NewVelty creates a terminator byte matcher
func NewVelty(inclusive bool, values ...byte) *veltyStart {
	return &veltyStart{values: values, inclusive: inclusive}
}
