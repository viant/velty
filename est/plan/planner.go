package plan

import (
	"fmt"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/cache"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/plan/scope"
	"github.com/viant/xunsafe"
	"reflect"
	"strconv"
)

const (
	fieldSeparator = "___"
)

type (
	Planner struct {
		Scope
		bufferSize int
		count      *int
		transients *int
		*est.Control
		Type         scope.Type
		selectors    *est.Selectors
		types        map[string]reflect.Type
		indexCounter int
		cache        *cache.Cache
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

	sel := est.NewSelector(name, id, sType, nil)
	if err := p.addSelector(sel); err != nil {
		return err
	}

	return nil
}

func (p *Planner) SelectorExpr(selector *expr.Select) *est.Selector {
	sel := p.selectorByName(selector.ID)
	if sel == nil {
		return nil
	}

	if selector.X == nil {
		return sel
	}

	call := selector.X
	parentType := sel.Type

	selectorId := selector.ID

	wasPtr := false
	for call != nil {
		if parentType.Kind() == reflect.Ptr {
			wasPtr = true
			parentType = deref(parentType)
		}

		switch actual := call.(type) {
		case *expr.Select:
			selectorId = selectorId + fieldSeparator + actual.ID
			field := xunsafe.FieldByName(parentType, actual.ID)
			if field == nil {
				return nil
			}

			sel = p.ensureSelector(selectorId, field, sel)
			sel.Indirect = wasPtr
			parentType = field.Type
			call = actual.X
		}
	}

	return sel
}

func deref(rType reflect.Type) reflect.Type {
	if rType.Kind() == reflect.Ptr {
		return deref(rType.Elem())
	}
	return rType
}

func (p *Planner) selectorID(name string) string {
	return p.prefix + name
}

func (p *Planner) accumulator(t reflect.Type) *est.Selector {
	name := "_T" + strconv.Itoa(*p.transients)
	*p.transients++
	sel := est.NewSelector(p.selectorID(name), name, t, nil)
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

	expr.Selector.Indirect = t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice
	return p.addSelector(expr.Selector)
}

func (p *Planner) addSelector(sel *est.Selector) error {
	if sel.ID == "" {
		return fmt.Errorf("selector ID was empty")
	}
	if sel.Type == nil {
		return fmt.Errorf("selector %v type was empty", sel.Name)
	}
	if p.selectorByName(sel.ID) != nil {
		return fmt.Errorf("variable %v already defined", sel.Name)
	}

	p.selectors.Append(sel)
	if sel.Parent == nil {
		sel.Field = p.Type.AddField(sel.Name, sel.Type)
	}
	return nil
}

func (p *Planner) selectorByName(name string) *est.Selector {
	if idx, ok := p.selectors.Index[p.prefix+name]; ok {
		return p.selectors.Selector(idx)
	}

	//check selector in upstream scopes
	for i := len(p.upstream) - 1; i >= 0; i-- {
		if idx, ok := p.selectors.Index[p.upstream[i]+name]; ok {
			return p.selectors.Selector(idx)
		}
	}

	//check global scope
	if idx, ok := p.selectors.Index[name]; ok {
		return p.selectors.Selector(idx)
	}
	return nil
}

func (p *Planner) compileIndexSelectorExpr() (*op.Expression, error) {
	p.indexCounter++
	fieldName := "_indexP_" + strconv.Itoa(p.indexCounter)

	indexType := reflect.TypeOf(0)
	selector := est.NewSelector(p.prefix+fieldName, fieldName, indexType, nil)
	if err := p.addSelector(selector); err != nil {
		return nil, err
	}

	return &op.Expression{
		Type:     indexType,
		Selector: selector,
	}, nil
}

func (p *Planner) ensureSelector(id string, field *xunsafe.Field, sel *est.Selector) *est.Selector {
	if selIndex, ok := p.selectors.Index[id]; ok {
		return p.selectors.Selector(selIndex)
	}

	selector := est.SelectorWithField(id, field, sel)
	if err := p.addSelector(selector); err != nil {
		return nil
	}

	return selector
}

func New(bufferSize int) *Planner {
	count := 0
	transients := 0
	ctl := est.Control(0)
	return &Planner{
		bufferSize: bufferSize,
		count:      &count,
		transients: &transients,
		Control:    &ctl,
		Type:       scope.Type{},
		Scope:      Scope{},
		selectors:  est.NewSelectors(),
		types:      map[string]reflect.Type{},
		cache:      cache.NewCache(),
	}
}
