package est

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

func (s *Selectors) Append(sel *Selector) {
	s.selectors = append(s.selectors, sel)
	s.Index[sel.ID] = len(s.selectors) - 1
}

func (s *Selectors) Merge(other *Selectors) {
	for i, _ := range other.selectors {
		s.Append(other.Selector(i))
	}
}

func (s *Selectors) Copy() *Selectors {
	result := NewSelectors()
	for i := range s.selectors {
		selector := s.selectors[i]
		result.Append(selector)
	}
	return result
}

func (s *Selectors) Selectors() []*Selector {
	return s.selectors
}
