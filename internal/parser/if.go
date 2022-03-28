package parser

import (
	"github.com/viant/parsly"
	astmt "github.com/viant/velty/internal/ast/stmt"
)

func matchIf(cursor *parsly.Cursor) (*astmt.If, error) {
	expression, err := matchEquationExpression(cursor)
	if err != nil {
		return nil, err
	}

	return &astmt.If{
		Condition: expression,
		Body:      astmt.Block{},
		Else:      nil,
	}, nil
}
