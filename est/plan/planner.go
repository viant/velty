package plan

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/plan/scope"
	"reflect"
	"strconv"
)

type (
	Planner struct {
		count      *int
		transients *int
		*est.Control
		Type scope.Type
		Scope
		selectors *[]*est.Selector
		index     map[string]int
		types     map[string]reflect.Type
	}
	Scope struct {
		upstream []string
		prefix   string
	}
)

func (s *Planner) AddNamedType(name string, t reflect.Type) {
	s.types[name] = t
}

func (s *Planner) AddType(t reflect.Type) {
	s.types[t.Name()] = t
}

func (s *Planner) NewScope() *Planner {
	*s.count++
	var types = make(map[string]reflect.Type)
	for k, v := range s.types {
		types[k] = v
	}
	return &Planner{
		Control: s.Control,
		Scope: Scope{
			upstream: append(s.upstream, s.prefix),
			prefix:   "s" + strconv.Itoa(*s.count),
		},
		selectors:  s.selectors,
		types:      types,
		transients: s.transients,
		count:      s.count,
	}
}

func (s *Planner) DefineVariable(name string, v interface{}) error {
	var sType reflect.Type
	switch t := v.(type) {
	case reflect.Type:
		sType = t
	default:
		t = reflect.TypeOf(v)
	}
	sel := est.NewSelector(name, name, sType)
	return s.addSelector(sel)
}

func (s *Planner) Selector(ID string) *est.Selector {
	if idx, ok := s.index[s.prefix+ID]; ok {
		return (*s.selectors)[idx]
	}
	//check selector in upstream scopes
	for i := len(s.upstream) - 1; i >= 0; i-- {
		if idx, ok := s.index[s.upstream[i]+ID]; ok {
			return (*s.selectors)[idx]
		}
	}
	//check global scope
	if idx, ok := s.index[ID]; ok {
		return (*s.selectors)[idx]
	}
	return nil
}

func (s *Planner) selectorID(name string) string {
	return s.prefix + name
}

func (s *Planner) AddRegistry(t reflect.Type) *est.Selector {
	name := "_T" + strconv.Itoa(*s.transients)
	*s.transients++
	sel := est.NewSelector(s.selectorID(name), name, t)
	_ = s.addSelector(sel)
	return sel
}

func (s *Planner) addSelector(sel *est.Selector) error {
	index := len(*s.selectors)
	if sel.ID == "" {
		return fmt.Errorf("selector ID was empty")
	}
	if sel.Type == nil {
		return fmt.Errorf("selector %v type was empty", sel.Name)
	}
	if s.Selector(sel.ID) != nil {
		return fmt.Errorf("variable %v already defined", sel.Name)
	}
	s.index[sel.ID] = index
	if len(sel.Ancestors) == 0 {
		sel.Field = s.Type.AddField(sel.Name, sel.Type)
	}
	*s.selectors = append(*s.selectors, sel)
	return nil
}

func New() *Planner {
	count := 0
	transients := 0
	ctl := est.Control(0)
	var selectors []*est.Selector
	return &Planner{
		count:      &count,
		transients: &transients,
		Control:    &ctl,
		Type:       scope.Type{},
		Scope:      Scope{},
		selectors:  &selectors,
		index:      map[string]int{},
		types:      map[string]reflect.Type{},
	}
}
