package parser

import (
	"fmt"
	"github.com/viant/parsly"
	"github.com/viant/velty/internal/ast/stmt"
)

func matchAssign(cursor *parsly.Cursor) (*stmt.Statement, error) {
	variable, err := matchVariable(cursor)
	if err != nil {
		return nil, err
	}

	tokenCandidates := []*parsly.Token{Assign}
	matched := cursor.MatchAfterOptional(WhiteSpace, tokenCandidates...)
	if matched.Code == parsly.EOF || matched.Code == parsly.Invalid {
		return nil, cursor.NewError(tokenCandidates...)
	}

	token := matchToken(matched)
	if token == "" {
		return nil, fmt.Errorf("didn't found operator token for given token %v", matched.Name)
	}

	_, expression, err := matchOperand(cursor, Boolean, String, Number)
	if err != nil {
		return nil, err
	}

	return &stmt.Statement{
		X:  variable,
		Op: token,
		Y:  expression,
	}, nil
}
