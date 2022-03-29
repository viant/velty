package parser

import (
	"fmt"
	"github.com/viant/velty/internal/ast"
	astmt "github.com/viant/velty/internal/ast/stmt"
)

type Builder struct {
	buffer []ast.StatementContainer
	block  astmt.Block
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (s *Builder) PushStatement(matchToken int, statement ast.Statement) error {
	switch matchToken {
	case elseIfToken, elseToken:
		lastNode := s.Last()
		if err := addIfExpression(lastNode, statement); err != nil {
			return err
		}
		return nil
	case endToken:
		if err := s.terminateStatement(); err != nil {
			return err
		}
		return nil
	}

	switch actual := statement.(type) {
	case ast.StatementContainer:
		if last := s.Last(); last != nil {
			last.AddStatement(actual)
		}
		s.buffer = append(s.buffer, actual)
	default:
		s.appendStatement(actual)
	}

	return nil
}

func addIfExpression(node ast.Node, expression ast.Node) error {
	switch nodeType := node.(type) {
	case astmt.ConditionContainer:
		switch exprType := expression.(type) {
		case *astmt.If:
			nodeType.AddCondition(exprType)
			return nil
		default:
			return fmt.Errorf("expected stmt.If but got %T", expression)
		}
	}
	return fmt.Errorf("expected stmt.Condition but got %T", node)
}

func (s *Builder) terminateStatement() error {
	if len(s.buffer) == 0 {
		return fmt.Errorf("unexpected expression end")
	}

	node := s.buffer[len(s.buffer)-1]
	s.buffer = s.buffer[:len(s.buffer)-1]

	if len(s.buffer) == 0 {
		s.block.AddStatement(node)
	}

	return nil
}

func (s *Builder) Last() ast.StatementContainer {
	if len(s.buffer) == 0 {
		return nil
	}
	return s.buffer[len(s.buffer)-1]
}

func (s *Builder) BufferSize() int {
	return len(s.buffer)
}

func (s *Builder) Block() *astmt.Block {
	return &s.block
}

func (s *Builder) appendStatement(statement ast.Statement) {
	if len(s.buffer) != 0 {
		s.Last().AddStatement(statement)
		return
	}

	s.block.AddStatement(statement)
}
