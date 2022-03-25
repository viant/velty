package op

import (
	"fmt"
	"github.com/viant/velty/utils"
)

type Selectors struct {
	selectors []*Selector
	Index     map[string]int
}

func NewSelectors() *Selectors {
	return &Selectors{
		selectors: make([]*Selector, 0),
		Index:     map[string]int{},
	}
}

func (s *Selectors) Selector(index int) *Selector {
	return (s.selectors)[index]
}

func (s *Selectors) Append(sel *Selector) error {
	if _, ok := s.Index[sel.ID]; ok {
		return fmt.Errorf("selector with id %v, already exists", sel.ID)
	}

	s.selectors = append(s.selectors, sel)
	s.Index[sel.ID] = len(s.selectors) - 1
	return nil
}

func (s *Selectors) Selectors() []*Selector {
	return s.selectors
}

func (s *Selectors) ById(selectorId string) (*Selector, bool) {
	selectorId = utils.UpperCaseFirstLetter(selectorId)

	index, found := s.Index[selectorId]
	if !found {
		return nil, found
	}

	return s.selectors[index], true
}
