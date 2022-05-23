package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/ast/stmt"
)

func matchIf(cursor *parsly.Cursor) (*stmt.If, error) {
	expression, err := matchEquationExpression(cursor)

	if err != nil {
		return nil, err
	}

	return &stmt.If{
		Condition: expression,
		Body:      stmt.Block{},
		Else:      nil,
	}, nil
}
