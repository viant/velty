package parser

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/stmt"
)

type Builder struct {
	buffer []ast.StatementContainer
	block  stmt.Block
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
		s.buffer = append(s.buffer, actual)
	default:
		s.appendStatement(actual)
	}

	return nil
}

func (s *Builder) terminateStatement() error {
	if len(s.buffer) == 0 {
		return fmt.Errorf("unexpected expression end")
	}

	node := s.buffer[len(s.buffer)-1]
	s.buffer = s.buffer[:len(s.buffer)-1]

	s.block.AddStatement(node)
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

func (s *Builder) Block() *stmt.Block {
	return &s.block
}

func (s *Builder) appendStatement(statement ast.Statement) {
	if len(s.buffer) != 0 {
		s.Last().AddStatement(statement)
		return
	}

	s.block.AddStatement(statement)
}
