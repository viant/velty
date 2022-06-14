package velty

import (
	"fmt"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/functions"
	"github.com/viant/velty/utils"
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
		Type      *est.Type
		selectors *op.Selectors
		constants *constants
		*op.Functions
		cache      *cache
		escapeHTML bool
	}
)

//EmbedVariable enrich the Type by adding Anonymous field with given name.
//val can be either of the reflect.Type or regular type (i.e. Foo)
func (p *Planner) EmbedVariable(val interface{}) error {
	var rType reflect.Type
	switch actual := val.(type) {
	case reflect.Type:
		rType = actual
	default:
		rType = reflect.TypeOf(val)
	}

	field := p.Type.EmbedType(rType)
	vTag := Parse(field.Tag.Get(velty))

	return p.addSelectors(vTag.Prefix, field)
}

func (p *Planner) addSelectors(prefix string, field reflect.StructField) error {
	return p.createSelectors(prefix, field, nil, 0, 0, false)
}

func (p *Planner) createSelectors(prefix string, field reflect.StructField, parent *op.Selector, offsetSoFar, initialOffset uintptr, indirect bool) error {
	if field.Anonymous {
		initialOffset += field.Offset
	}

	indirect = indirect || field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Slice
	vTag := Parse(field.Tag.Get(velty))
	if err := p.indexSelectorIfNeeded(prefix, field, parent, vTag, offsetSoFar, initialOffset, indirect); err != nil {
		return err
	}

	if !field.Anonymous {
		offsetSoFar += field.Offset
		initialOffset = 0
	}

	rType, _ := dereference(field)
	if rType.Kind() == reflect.Struct {
		for i := 0; i < rType.NumField(); i++ {
			actualParent := p.ensureStructSelector(field, prefix)

			childPrefix := vTag.Prefix
			if !field.Anonymous {
				childPrefix = field.Name + fieldSeparator
			}

			err := p.createSelectors(prefix+childPrefix, rType.Field(i), actualParent, offsetSoFar, initialOffset, indirect)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Planner) indexSelectorIfNeeded(prefix string, field reflect.StructField, parent *op.Selector, vTag *Tag, offset uintptr, anonymousOffset uintptr, indirect bool) error {
	if field.Anonymous {
		return nil
	}

	fieldNames := []string{field.Name}
	if len(vTag.Names) != 0 {
		fieldNames = vTag.Names
	}

	newField := xunsafe.NewField(field)
	newField.Offset += anonymousOffset
	valueField := xunsafe.NewField(reflect.StructField{Name: "TEMP", Offset: 0, Type: field.Type})
	var err error
	for _, name := range fieldNames {
		fieldSelector := op.SelectorWithField(prefix+name, newField, parent, indirect, offset)
		if err = p.selectors.Append(fieldSelector); err != nil {
			return fmt.Errorf("%w, prefix is required, if parent field is an Anonymous, and any other parent field has the same name", err)
		}
		fieldSelector.ValueField = valueField
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

func (p *Planner) ensureStructSelector(field reflect.StructField, prefix string) *op.Selector {
	sel, _ := p.selectors.ById(prefix + field.Name)
	return sel
}

//DefineVariable enrich the Type by adding field with given name.
//val can be either of the reflect.Type or regular type (i.e. Foo)
func (p *Planner) DefineVariable(name string, v interface{}) error {
	if p.selectorByName(name) != nil {
		return nil
	}

	name = utils.UpperCaseFirstLetter(name)

	var sType reflect.Type
	switch t := v.(type) {
	case reflect.Type:
		sType = t
	default:
		sType = reflect.TypeOf(v)
	}

	field := p.Type.AddField(name, name, sType)
	return p.addSelectors("", field)
}

func (p *Planner) selector(selector *expr.Select) (*op.Selector, error) {
	resultSelector := p.selectorByName(selector.ID)
	if resultSelector == nil {
		return nil, nil
	}

	if selector.X == nil {
		return resultSelector, nil
	}

	call := selector.X

	parentType := resultSelector.Type
	selectorId := selector.ID

	for call != nil {
		parentType = deref(parentType)
		switch actual := call.(type) {
		case *expr.Select:
			selectorId = selectorId + fieldSeparator + actual.ID

			if actual.X != nil {
				switch next := actual.X.(type) {
				case *expr.Call:
					var err error
					resultSelector, err = p.newFuncSelector(selectorId, actual, next, resultSelector)
					if err != nil {
						return nil, err
					}

					parentType = resultSelector.Func.ResultType
					call = next.X
					continue
				}
			}

			field, err := p.fieldByName(parentType, actual, selectorId)
			if err != nil {
				return nil, err
			}

			var found bool
			resultSelector, found = p.selectors.ById(selectorId)
			if !found {
				return nil, fmt.Errorf("not found selector for the %v", strings.ReplaceAll(selectorId, fieldSeparator, "."))
			}

			parentType = field.Type
			call = actual.X
		}
	}

	return resultSelector, nil
}

func (p *Planner) fieldByName(parentType reflect.Type, actual *expr.Select, selectorId string) (*xunsafe.Field, error) {
	field := xunsafe.FieldByName(parentType, actual.ID)
	if field != nil {
		if Parse(field.Tag.Get(velty)).Omit {
			return nil, fmt.Errorf("can't create selector for field %v", field.Name)
		}
		return field, nil
	}

	for i := 0; i < parentType.NumField(); i++ {
		vTag := Parse(parentType.Field(i).Tag.Get(velty))
		if vTag.nameEqual(actual.ID) {
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

func (p *Planner) accumulator(t reflect.Type) *op.Selector {
	name := "_T" + strconv.Itoa(*p.transients)
	*p.transients++
	sel := op.NewSelector(name, name, t, nil)
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

	if err := p.DefineVariable(expr.Name, t); err != nil {
		return err
	}

	expr.Type = t
	expr.Selector = p.selectorByName(expr.Name)
	return nil
}

func (p *Planner) validateSelector(sel *op.Selector) error {
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

func (p *Planner) selectorByName(name string) *op.Selector {
	if idx, ok := p.selectors.Index[name]; ok {
		return p.selectors.Selector(idx)
	}
	return nil
}

func (p *Planner) newFuncSelector(selectorId string, field *expr.Select, call *expr.Call, prev *op.Selector) (*op.Selector, error) {
	var err error
	aFunc, ok := p.Functions.Method(prev.Type, field.ID)
	if !ok {
		return nil, fmt.Errorf("not found function: %v", field.ID)
	}

	name := "_T" + strconv.Itoa(*p.transients)
	*p.transients++
	strField := p.Type.AddField(name, name, aFunc.ResultType)

	operands, err := p.selectorOperands(call, prev)
	if err != nil {
		return nil, err
	}

	newSelector := op.FunctionSelector(selectorId, strField, aFunc, call.Args, prev)
	newSelector.Args = operands
	newSelector.Type = aFunc.ResultType

	return newSelector, nil
}

func (p *Planner) selectorOperands(call *expr.Call, prev *op.Selector) ([]*op.Operand, error) {
	var err error
	operands := make([]*op.Operand, len(call.Args)+1)
	operands[0], err = op.NewExpression(prev).Operand(*p.Control)

	if err != nil {
		return nil, err
	}

	for i := 1; i < len(operands); i++ {
		expression, err := p.compileExpr(call.Args[i-1])
		if err != nil {
			return nil, err
		}

		operand, err := expression.Operand(*p.Control)
		if err != nil {
			return nil, err
		}
		operands[i] = operand
	}
	return operands, nil
}

func New(options ...Option) *Planner {
	transients := 0
	ctl := est.Control(0)
	planner := &Planner{
		transients: &transients,
		Control:    &ctl,
		Type:       est.NewType(),
		selectors:  op.NewSelectors(),
		cache:      newCache(0),
		Functions:  op.NewFunctions(),
		constants:  newConstants(),
	}

	planner.init(options)

	return planner
}

func (p *Planner) New() *Planner {
	scope := &Planner{
		bufferSize: p.bufferSize,
		transients: p.transients,
		Control:    p.Control,
		Type:       p.Type.Snapshot(),
		selectors:  p.selectors.Snapshot(),
		constants:  p.constants,
		Functions:  p.Functions,
		cache:      p.cache,
		escapeHTML: p.escapeHTML,
	}

	return scope
}

func (p *Planner) apply(options []Option) {
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			p.bufferSize = int(actual)
		case CacheSize:
			p.cache = newCache(int(actual))
		case EscapeHTML:
			p.escapeHTML = bool(actual)
		}
	}
}

func (p *Planner) registerConst(i *[]int) {
	p.constants.add(i)
}

func (p *Planner) init(options []Option) {
	p.apply(options)

	_ = p.DefineVariable("strings", functions.Strings{})
	_ = p.DefineVariable("math", functions.Math{})
	_ = p.DefineVariable("strconv", functions.Strconv{})
	_ = p.DefineVariable("slices", functions.Slices{})
	_ = p.DefineVariable("types", functions.Types{})
	_ = p.DefineVariable("errors", functions.Errors{})
}
