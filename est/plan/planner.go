package plan

import (
	"fmt"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/est/plan/scope"
	"github.com/viant/velty/tag"
	"github.com/viant/xunsafe"
	"reflect"
	"strconv"
	"strings"
)

const (
	fieldSeparator = "___"
)

type (
	Planner struct {
		bufferSize int
		transients *int
		*est.Control
		Type      *scope.Type
		selectors *est.Selectors
		cache     *Cache
	}
)

func (p *Planner) EmbedType(name string, val interface{}) error {
	var rType reflect.Type
	switch actual := val.(type) {
	case reflect.Type:
		rType = actual
	default:
		rType = reflect.TypeOf(val)
	}

	field := p.Type.EmbedType(name, name, rType)
	vTag := tag.Parse(field.Tag.Get(tag.Velty))

	return p.createSelectors(vTag.Prefix, field, nil)
}

func (p *Planner) createSelectors(prefix string, field reflect.StructField, parent *est.Selector) error {
	var err error

	vTag := tag.Parse(field.Tag.Get(tag.Velty))

	if !field.Anonymous {
		fieldNames := []string{field.Name}
		if len(vTag.Names) != 0 {
			fieldNames = vTag.Names
		}

		for _, name := range fieldNames {
			fieldSelector := est.SelectorWithField(prefix+name, xunsafe.NewField(field), parent)
			parent = fieldSelector
			if err = p.selectors.Append(fieldSelector); err != nil {
				return fmt.Errorf("%w, you have to specify prefix, if parent field is an Anonymous, and any other parent field has the same name", err)
			}
		}
	}

	rType, wasPtr := dereference(field)
	if rType.Kind() == reflect.Struct {
		for i := 0; i < rType.NumField(); i++ {

			actualParent := parent
			if wasPtr {
				actualParent = p.ensureStructSelector(field, prefix)
			}

			childPrefix := vTag.Prefix
			if !field.Anonymous {
				childPrefix = field.Name + fieldSeparator
			}

			err = p.createSelectors(prefix+childPrefix, rType.Field(i), actualParent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dereference(field reflect.StructField) (reflect.Type, bool) {
	rType := field.Type
	wasPtr := false
	if rType.Kind() == reflect.Ptr {
		wasPtr = true
		rType = rType.Elem()
	}
	return rType, wasPtr
}

func (p *Planner) ensureStructSelector(field reflect.StructField, prefix string) *est.Selector {
	sel, _ := p.selectors.ById(prefix + field.Name)
	return sel
}

func (p *Planner) DefineVariable(name string, v interface{}) error {
	var sType reflect.Type
	switch t := v.(type) {
	case reflect.Type:
		sType = t
	default:
		sType = reflect.TypeOf(v)
	}

	field := p.Type.AddField(name, name, sType)
	return p.createSelectors("", field, nil)
}

func (p *Planner) SelectorExpr(selector *expr.Select) (*est.Selector, error) {
	sel := p.selectorByName(selector.ID)
	if sel == nil {
		return nil, nil
	}

	if selector.X == nil {
		return sel, nil
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
			field, err := p.fieldByName(parentType, actual, selectorId)
			if err != nil {
				return nil, err
			}

			var found bool
			sel, found = p.selectors.ById(selectorId)
			if !found {
				return nil, fmt.Errorf("not found selector for the %v", strings.ReplaceAll(selectorId, fieldSeparator, "."))
			}
			sel.Indirect = wasPtr
			parentType = field.Type
			call = actual.X
		}
	}

	return sel, nil
}

func (p *Planner) fieldByName(parentType reflect.Type, actual *expr.Select, selectorId string) (*xunsafe.Field, error) {
	field := xunsafe.FieldByName(parentType, actual.ID)
	if field != nil {
		return field, nil
	}

	for i := 0; i < parentType.NumField(); i++ {
		vTag := tag.Parse(parentType.Field(i).Tag.Get(tag.Velty))
		if vTag.NameEqual(actual.ID) {
			return xunsafe.NewField(parentType.Field(i)), nil
		}
	}

	return nil, fmt.Errorf("not found field %v at %v", strings.ReplaceAll(selectorId, fieldSeparator, "."), parentType.String())
}

func deref(rType reflect.Type) reflect.Type {
	if rType.Kind() == reflect.Ptr {
		return deref(rType.Elem())
	}
	return rType
}

func (p *Planner) accumulator(t reflect.Type) *est.Selector {
	name := "_T" + strconv.Itoa(*p.transients)
	*p.transients++
	sel := est.NewSelector(name, name, t, nil)
	if t != nil {
		_ = p.selectors.Append(sel)
		sel.Field = xunsafe.NewField(p.Type.AddField(name, name, t))
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

	if err := p.validateSelector(expr.Selector); err != nil {
		return err
	}

	field := p.Type.AddField(expr.ID, strings.Title(expr.Name), t)
	expr.Field = xunsafe.NewField(field)
	return p.selectors.Append(expr.Selector)
}

func (p *Planner) validateSelector(sel *est.Selector) error {
	if sel.ID == "" {
		return fmt.Errorf("selector ID was empty")
	}

	if sel.Type == nil {
		return fmt.Errorf("selector %v type was empty", sel.Name)
	}

	if p.selectorByName(sel.ID) != nil {
		return fmt.Errorf("variable %v already defined", sel.Name)
	}

	return nil
}

func (p *Planner) selectorByName(name string) *est.Selector {
	if idx, ok := p.selectors.Index[name]; ok {
		return p.selectors.Selector(idx)
	}
	return nil
}

func New(sizes ...int) *Planner {
	bufferSize := 0
	if len(sizes) > 0 {
		bufferSize = sizes[0]
	}

	cacheSize := 0
	if len(sizes) > 1 {
		cacheSize = sizes[1]
	}

	transients := 0
	ctl := est.Control(0)
	result := &Planner{
		bufferSize: bufferSize,
		transients: &transients,
		Control:    &ctl,
		Type:       scope.NewType(),
		selectors:  est.NewSelectors(),
		cache:      NewCache(cacheSize),
	}

	return result
}
