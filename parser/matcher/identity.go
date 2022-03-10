package matcher

import (
	"github.com/viant/parsly"
)

type identity struct {
	fullMatch   bool
	onlyLetters bool
}

//Match matches a string
func (n *identity) Match(cursor *parsly.Cursor) (matched int) {
	input := cursor.Input
	pos := cursor.Pos
	if startsWithCharacter := IsLetter(input[pos]); startsWithCharacter {
		pos++
		matched++
	} else {
		return 0
	}

	size := len(input)
	for i := pos; i < size; i++ {
		switch input[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_', '.':
			matched++
			continue
		case '\n', '\r', ' ':
			if n.fullMatch {
				return 0
			}
			return matched
		case '(':
			return matched
		default:
			if IsLetter(input[i]) {
				matched++
				continue
			}

			if n.fullMatch {
				return 0
			} else {
				return matched
			}
		}
	}

	return matched
}

func IsLetter(b byte) bool {
	if (b < 'a' || b > 'z') && (b < 'A' || b > 'Z') {
		return false
	}
	return true
}

//NewIdentity creates a string matcher
func NewIdentity(fullMatch bool, onlyLetters bool) *identity {
	return &identity{
		fullMatch:   fullMatch,
		onlyLetters: onlyLetters,
	}
}
