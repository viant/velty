package parser

import "github.com/viant/parsly"

func isUnaryMatched(tokenMatch *parsly.TokenMatch) bool {
	switch tokenMatch.Code {
	case negationToken:
		return true
	}
	return false
}
