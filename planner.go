package velty

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/xunsafe"
	"reflect"
	"strings"
)

const (
	fieldSeparator = "___"
)

type (
	Planner struct {
		bufferSize int
		*est.Control
		Type      *est.Type
		selectors *op.Selectors
		constants *constants
		*op.Functions
		cache        *cache
		escapeHTML   bool
		panicOnError bool
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

	return p.addSelectors(vTag.Prefix, field, field.Name)
}

func (p *Planner) addSelectors(prefix string, field reflect.StructField, fieldName string) error {
	detector := NewCycleDetector(field.Type)
	return p.createSelectors(prefix, field, nil, 0, 0, false, detector, fieldName)
}

func (p *Planner) createSelectors(prefix string, field reflect.StructField, parent *op.Selector, offsetSoFar, initialOffset uintptr, indirect bool, cycleDetector *CycleDetector, fieldName string) error {
	cycleNode, cycleSelector := p.cycle(cycleDetector, field, parent)

	if field.Anonymous {
		initialOffset += field.Offset
	}

	indirect = indirect || field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Slice
	vTag := Parse(field.Tag.Get(velty))

	parent, err := p.indexSelectorIfNeeded(prefix, field, parent, offsetSoFar, initialOffset, indirect, cycleSelector, fieldName)
	if err != nil || cycleSelector != nil {
		return err
	}

	if !field.Anonymous {
		offsetSoFar += field.Offset
		initialOffset = 0

		if prefix == "" {
			prefix = fieldName + fieldSeparator
		} else {
			prefix = prefix + fieldName + fieldSeparator
		}
	} else {
		if vTag.Prefix != "" {
			prefix += vTag.Prefix
		}
	}

	return p.addChildrenSelectors(prefix, field, offsetSoFar, initialOffset, indirect, cycleNode, parent)
}

func (p *Planner) cycle(cycleDetector *CycleDetector, field reflect.StructField, parent *op.Selector) (*CycleDetector, *op.Selector) {
	child, cycle := cycleDetector.Child(field.Type, parent)
	if cycle {
		return child, child.parentSelector
	}
	return child, nil
}

func (p *Planner) addChildrenSelectors(holderPrefix string, field reflect.StructField, offsetSoFar, initialOffset uintptr, indirect bool, detector *CycleDetector, parent *op.Selector) error {
	rType, elemed := elemIfNeeded(field.Type)
	if rType.Kind() == reflect.Struct {
		for i := 0; i < rType.NumField(); i++ {

			structField := rType.Field(i)
			vTag := Parse(structField.Tag.Get(velty))
			fieldNames := []string{structField.Name}
			if len(vTag.Names) != 0 {
				fieldNames = vTag.Names
			}

			for _, name := range fieldNames {
				if err := p.createSelectors(holderPrefix, structField, parent, offsetSoFar, initialOffset, indirect || elemed, detector, name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Planner) indexSelectorIfNeeded(prefix string, field reflect.StructField, parent *op.Selector, offset uintptr, anonymousOffset uintptr, indirect bool, cycleSelector *op.Selector, name string) (*op.Selector, error) {
	if field.Anonymous && field.Type.Kind() != reflect.Ptr {
		return parent, nil
	}

	newField := xunsafe.NewField(field)
	newField.Offset += anonymousOffset
	var err error
	var fieldSelector *op.Selector
	if cycleSelector != nil {
		fieldSelector = op.NewCycleSelector(prefix+name, newField, parent, indirect, offset, cycleSelector)
	} else {
		fieldSelector = op.SelectorWithField(prefix+name, newField, parent, indirect, offset)
	}

	if field.Anonymous {
		return fieldSelector, nil
	}

	if err = p.selectors.Append(fieldSelector); err != nil {
		return nil, fmt.Errorf("%w, prefix is required, if parent field is an Anonymous, and any other parent field has the same name", err)
	}

	return fieldSelector, nil
}

func elemIfNeeded(rType reflect.Type) (reflect.Type, bool) {
	wasPtr := false
	for rType.Kind() == reflect.Ptr || rType.Kind() == reflect.Slice || rType.Kind() == reflect.Map {
		wasPtr = true
		rType = rType.Elem()
	}

	return rType, wasPtr
}

//DefineVariable enrich the Type by adding field with given name.
//val can be either of the reflect.Type or regular type (i.e. Foo)
func (p *Planner) DefineVariable(name string, v interface{}, names ...string) error {
	if p.selectorByName(name) != nil {
		return nil
	}

	var sType reflect.Type
	switch t := v.(type) {
	case reflect.Type:
		sType = t
	default:
		sType = reflect.TypeOf(v)
	}

	field := p.Type.AddField(name, name, sType)
	if err := p.addSelectors("", field, name); err != nil {
		return err
	}

	for _, additionalFieldName := range names {
		if err := p.addSelectors("", field, additionalFieldName); err != nil {
			return err
		}
	}

	return nil
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

	upstreamSelector := p.copyWithParent(resultSelector, resultSelector.Parent)
	var err error
	for call != nil {
		parentType = deref(parentType)
		resultSelector, call, err = p.matchSelector(call, resultSelector, selectorId, parentType)
		if err != nil {
			return nil, err
		}

		parentType = resultSelector.Type
		selectorId = resultSelector.ID
		upstreamSelector = p.copyWithParent(resultSelector, upstreamSelector)
	}

	return upstreamSelector, nil
}

func (p *Planner) copyWithParent(dest, parent *op.Selector) *op.Selector {
	selCopy := *dest
	selCopy.Parent = parent
	return &selCopy
}

func (p *Planner) matchSelector(call ast.Expression, resultSelector *op.Selector, selectorId string, parentType reflect.Type) (*op.Selector, ast.Expression, error) {
	selector, next, err := p.tryMatchCall(call, resultSelector, selectorId)
	if err != nil || selector != nil {
		return selector, next, err
	}

	switch actual := call.(type) {
	case *expr.Select:
		if callSelector, callNext, callErr := p.tryMatchCall(actual.X, resultSelector, actual.ID); callSelector != nil || callErr != nil {
			return callSelector, callNext, callErr
		}

		_, err = p.fieldByName(parentType, actual, actual.ID)
		if err != nil {
			return nil, nil, err
		}

		selectorId = selectorId + fieldSeparator + actual.ID
		var found bool
		resultSelector, found = p.selectors.ById(selectorId)
		if !found {
			return nil, nil, fmt.Errorf("not found selector for the %v", strings.ReplaceAll(selectorId, fieldSeparator, "."))
		}

		return resultSelector, actual.X, nil
	}

	return resultSelector, nil, nil
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
	for {
		switch rType.Kind() {
		case reflect.Ptr, reflect.Slice:
			rType = rType.Elem()
		default:
			return rType
		}
	}
}

func (p *Planner) accumulator(t reflect.Type) *op.Selector {
	name := p.newName()
	sel := op.NewSelector(name, name, t, nil)
	if t != nil {
		_ = p.selectors.Append(sel)
		sel.Field = xunsafe.NewField(p.Type.AddFieldWithTag(name, name, "", t))
	}
	return sel
}

func (p *Planner) newName() string {
	return p.Type.ReserveNewName()
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

	if funcSelector, ok := p.Functions.FuncSelector(name, nil); ok {
		return funcSelector
	}

	return nil
}

func (p *Planner) newFuncSelector(selectorId string, methodName string, call *expr.Call, prev *op.Selector) (*op.Selector, error) {
	var err error
	aFunc, err := p.Functions.Method(prev.Type, methodName, call)
	if err != nil {
		return nil, fmt.Errorf("not found function %v, due to: %w", methodName, err)
	}

	operands, err := p.selectorOperands(call, prev)
	if err != nil {
		return nil, err
	}

	accumulator := p.accumulator(aFunc.ResultType)
	newSelector := op.FunctionSelector(selectorId, accumulator.Field, aFunc, prev)
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
	ctl := est.Control(0)

	planner := &Planner{
		Control:   &ctl,
		Type:      est.NewType(),
		selectors: op.NewSelectors(),
		cache:     newCache(0),
		constants: newConstants(),
	}

	planner.init(options)

	return planner
}

func (p *Planner) New() *Planner {
	scope := &Planner{
		bufferSize: p.bufferSize,
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

func (p *Planner) apply(options []Option) *op.Func {
	var aFunc *op.Func
	for _, option := range options {
		switch actual := option.(type) {
		case BufferSize:
			p.bufferSize = int(actual)
		case CacheSize:
			p.cache = newCache(int(actual))
		case EscapeHTML:
			p.escapeHTML = bool(actual)
		case PanicOnError:
			p.panicOnError = bool(actual)
		case *op.Func:
			aFunc = actual
		}
	}

	return aFunc
}

func (p *Planner) registerConst(i *[]int) {
	p.constants.add(i)
}

func (p *Planner) init(options []Option) {
	aFunc := p.apply(options)

	if aFunc == nil {
		p.Functions = op.NewFunctions()
	}
}

func (p *Planner) tryMatchCall(call ast.Expression, selector *op.Selector, ID string) (*op.Selector, ast.Expression, error) {
	if call == nil {
		return nil, nil, nil
	}

	matchCall, expression, err := p.matchCall(call, selector, ID)
	if err != nil {
		return nil, nil, err
	}

	return matchCall, expression, err
}

func (p *Planner) matchCall(call ast.Expression, selector *op.Selector, ID string) (*op.Selector, ast.Expression, error) {
	switch actual := call.(type) {
	case *expr.Call:
		callSelector, err := p.newFuncSelector(ID, ID, actual, selector)
		if err != nil {
			return nil, nil, err
		}

		return callSelector, actual.X, nil
	case *expr.SliceIndex:
		switch selector.Type.Kind() {
		case reflect.Map:
			mapSelector, err := p.newMapSelector(ID, actual, selector)
			if err != nil {
				return nil, nil, err
			}

			return mapSelector, actual.Y, nil
		case reflect.Interface:
			interfaceSelector, err := p.newInterfaceSelector(ID, actual, selector)
			if err != nil {
				return nil, nil, err
			}

			return interfaceSelector, actual.Y, nil

		default:
			sliceSelector, err := p.newSliceSelector(ID, actual, selector)
			if err != nil {
				return nil, nil, err
			}

			return sliceSelector, actual.Y, nil
		}
	}

	return nil, nil, nil
}

func (p *Planner) newSliceSelector(id string, actual *expr.SliceIndex, selector *op.Selector) (*op.Selector, error) {
	indexEpression, err := p.compileExpr(actual.X)
	if err != nil {
		return nil, err
	}

	operandExpression, err := indexEpression.Operand(*p.Control)
	if err != nil {
		return nil, err
	}

	sliceOperand, err := op.NewExpression(selector).Operand(*p.Control)
	if err != nil {
		return nil, err
	}

	return op.SliceSelector(id, "", sliceOperand, operandExpression, selector)
}

func (p *Planner) newMapSelector(id string, actual *expr.SliceIndex, selector *op.Selector) (*op.Selector, error) {
	keyOperand, err := p.compileOperand(actual.X)
	if err != nil {
		return nil, err
	}

	mapOperand, err := op.NewExpression(selector).Operand(*p.Control)
	if err != nil {
		return nil, err
	}

	return op.NewMapSelector(id, "", mapOperand, keyOperand, selector)
}

func (p *Planner) newInterfaceSelector(id string, actual *expr.SliceIndex, selector *op.Selector) (*op.Selector, error) {
	xOperand, err := op.NewExpression(selector).Operand(*p.Control)
	if err != nil {
		return nil, err
	}

	yOperand, err := p.compileOperand(actual.X)
	if err != nil {
		return nil, err
	}

	return op.NewInterfaceSelector(id, "", xOperand, yOperand, selector)
}

func (p *Planner) compileOperand(actual ast.Expression) (*op.Operand, error) {
	if actual == nil {
		return nil, nil
	}

	xExpr, err := p.compileExpr(actual)
	if err != nil {
		return nil, err
	}

	xOperand, err := xExpr.Operand(*p.Control)
	if err != nil {
		return nil, err
	}
	return xOperand, nil
}

func (p *Planner) derefHolderSelector(field reflect.StructField) *op.Selector {
	return &op.Selector{
		ID:           "",
		Type:         field.Type,
		Field:        xunsafe.NewField(field),
		Indirect:     true,
		ParentOffset: field.Offset,
	}
}
