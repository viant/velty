package plan

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/plan/scope"
	"reflect"
	"strconv"
)

const currentPkg = "github.com/viant/velty/est/plan"

type (
	Planner struct {
		bufferSize int
		count      *int
		transients *int
		*est.Control
		Type scope.Type
		Scope
		selectors    *[]*est.Selector
		index        map[string]int
		types        map[string]reflect.Type
		indexCounter int
	}
	Scope struct {
		upstream []string
		prefix   string
	}
)

func (p *Planner) AddNamedType(name string, t reflect.Type) {
	p.types[name] = t
}

func (p *Planner) AddType(t reflect.Type) {
	p.types[t.Name()] = t
}

func (p *Planner) NewScope() *Planner {
	*p.count++
	var types = make(map[string]reflect.Type)
	for k, v := range p.types {
		types[k] = v
	}
	return &Planner{
		Control: p.Control,
		Scope: Scope{
			upstream: append(p.upstream, p.prefix),
			prefix:   "p" + strconv.Itoa(*p.count),
		},
		selectors:  p.selectors,
		types:      types,
		transients: p.transients,
		count:      p.count,
	}
}

func (p *Planner) DefineVariable(name string, v interface{}) error {
	return p.defineVariable(name, name, v)
}

func (p *Planner) defineVariable(name string, id string, v interface{}) error {
	var sType reflect.Type
	switch t := v.(type) {
	case reflect.Type:
		sType = t
	default:
		sType = reflect.TypeOf(v)
	}

	sel := est.NewSelector(name, id, sType)
	return p.addSelector(sel)
}

func (p *Planner) Selector(name string) *est.Selector {

	//TODO add support for . in name

	if idx, ok := p.index[p.prefix+name]; ok {
		return (*p.selectors)[idx]
	}
	//check selector in upstream scopes
	for i := len(p.upstream) - 1; i >= 0; i-- {
		if idx, ok := p.index[p.upstream[i]+name]; ok {
			return (*p.selectors)[idx]
		}
	}
	//check global scope
	if idx, ok := p.index[name]; ok {
		return (*p.selectors)[idx]
	}
	return nil
}

func (p *Planner) selectorID(name string) string {
	return p.prefix + name
}

func (p *Planner) accumulator(t reflect.Type) *est.Selector {
	name := "_T" + strconv.Itoa(*p.transients)
	*p.transients++
	sel := est.NewSelector(p.selectorID(name), name, t)
	if t != nil {
		_ = p.addSelector(sel)
	}
	return sel
}

func (p *Planner) adjustSelector(expr *op.Expression, t reflect.Type) error {
	if expr.Selector.Type != nil {
		return nil
	}
	expr.Type = t
	expr.Selector.Type = t
	return p.addSelector(expr.Selector)
}

func (p *Planner) addSelector(sel *est.Selector) error {
	index := len(*p.selectors)
	if sel.ID == "" {
		return fmt.Errorf("selector ID was empty")
	}
	if sel.Type == nil {
		return fmt.Errorf("selector %v type was empty", sel.Name)
	}
	if p.Selector(sel.ID) != nil {
		return fmt.Errorf("variable %v already defined", sel.Name)
	}
	p.index[sel.ID] = index
	if len(sel.Ancestors) == 0 {
		sel.Field = p.Type.AddField(sel.Name, sel.Type)
	}
	*p.selectors = append(*p.selectors, sel)
	return nil
}

func (p *Planner) compileIndexSelectorExpr() (*op.Expression, error) {
	p.indexCounter++
	fieldName := "_indexP_" + strconv.Itoa(p.indexCounter)

	indexType := reflect.TypeOf(0)
	selector := est.NewSelector(p.prefix+fieldName, fieldName, indexType)
	if err := p.addSelector(selector); err != nil {
		return nil, err
	}

	return &op.Expression{
		Type:     indexType,
		Selector: selector,
	}, nil
}

func New(bufferSize int) *Planner {
	count := 0
	transients := 0
	ctl := est.Control(0)
	var selectors []*est.Selector
	return &Planner{
		bufferSize: bufferSize,
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
