package parser

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/stmt"
)

type Stack struct {
	Nodes []ast.Statement
	block stmt.Block
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(node ast.Statement) {
	s.Nodes = append(s.Nodes, node)
}

func (s *Stack) TransferToBlock() error {
	if len(s.Nodes) == 0 {
		return fmt.Errorf("unexpected expression end")
	}

	node := s.Nodes[len(s.Nodes)-1]
	s.Nodes = s.Nodes[:len(s.Nodes)-1]

	return s.block.AddStatement(node)
}

func (s *Stack) Last() ast.Statement {
	if len(s.Nodes) == 0 {
		return nil
	}
	return s.Nodes[len(s.Nodes)-1]
}

func (s *Stack) Size() int {
	return len(s.Nodes)
}

func (s *Stack) Block() *stmt.Block {
	return &s.block
}

func (s *Stack) AppendStatement(statement ast.Statement) error {
	if len(s.Nodes) != 0 {
		return s.Last().AddStatement(statement)
	}

	return s.block.AddStatement(statement)
}
