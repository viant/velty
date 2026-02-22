package parser

import (
	"github.com/viant/parsly"
	"github.com/viant/velty/ast/stmt"
)

func matchIf(cursor *parsly.Cursor, spans *spanState) (*stmt.If, error) {
	expression, err := matchEquationExpression(cursor, spans)

	if err != nil {
		return nil, err
	}

	return &stmt.If{
		Condition: expression,
		Body:      stmt.Block{},
		Else:      nil,
	}, nil
}
