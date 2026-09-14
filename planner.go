package velty

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/xreflect"
	"github.com/viant/xunsafe"
	"reflect"
	"time"
)

const (
	fieldSeparator = "___"
)

var TimeType = reflect.TypeOf(time.Time{})

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
		nodeSeq      int
		// listener receives parse events (reserved for future use)
		listener ParserListener
		// adjuster applies node transformations (reserved for future use)
		adjuster NodeAdjuster
		// evaluate traversal configuration
		evalConfig EvaluateConfig
		// planListener receives binding/type resolution hooks
		planListener PlannerListener
	}
)

// EmbedVariable enrich the Type by adding Anonymous field with given name.
// val can be either of the reflect.Type or regular type (i.e. Foo)
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
	var detector *CycleDetector
	return p.createSelectors(prefix, field, nil, 0, 0, false, detector, fieldName)
}

func (p *Planner) createSelectors(prefix string, field reflect.StructField, parent *op.Selector, offsetSoFar, initialOffset uintptr, indirect bool, cycleDetector *CycleDetector, fieldName string) error {
	cycleNode, hasCycle := cycleDetector.Child(field.Type, parent)

	if field.Anonymous {
		initialOffset += field.Offset
	}

	indirect = indirect || field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Slice
	vTag := Parse(field.Tag.Get(velty))

	newParent, err := p.indexSelectorIfNeeded(prefix, field, parent, offsetSoFar, initialOffset, indirect, fieldName)
	if err != nil || hasCycle {
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

	fieldType, _ := elemIfNeeded(field.Type)
	if fieldType == TimeType {
		return nil
	}

	return p.addChildrenSelectors(prefix, field, offsetSoFar, initialOffset, indirect, cycleNode, newParent)
}

func (p *Planner) addChildrenSelectors(holderPrefix string, field reflect.StructField, offsetSoFar, initialOffset uintptr, indirect bool, detector *CycleDetector, parent *op.Selector) error {
	rType, elemed := elemIfNeeded(field.Type)
	if rType.Kind() == reflect.Struct {
		fieldsLen := rType.NumField()
		for i := 0; i < fieldsLen; i++ {

			structField := rType.Field(i)
			vTag := Parse(structField.Tag.Get(velty))
			if vTag.Omit {
				continue
			}

			fieldNames := []string{structField.Name}
			if len(vTag.Names) != 0 {
				fieldNames = vTag.Names
			}

			for _, name := range fieldNames {
				if name == "_" {
					continue
				}

				if err := p.createSelectors(holderPrefix, structField, parent, offsetSoFar, initialOffset, indirect || elemed, detector, name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Planner) indexSelectorIfNeeded(prefix string, field reflect.StructField, parent *op.Selector, offset uintptr, anonymousOffset uintptr, indirect bool, name string) (*op.Selector, error) {
	if field.Anonymous && field.Type.Kind() != reflect.Ptr {
		return parent, nil
	}

	newField := xunsafe.NewField(field)
	newField.Offset += anonymousOffset
	var err error
	fieldSelector := op.SelectorWithField(prefix+name, newField, parent, indirect, offset)
	fieldSelector.NodeID = p.nextNodeID()

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
	seen := map[reflect.Type]bool{}
	for rType.Kind() == reflect.Ptr || rType.Kind() == reflect.Slice || rType.Kind() == reflect.Map {
		if seen[rType] {
			break
		}
		seen[rType] = true
		wasPtr = true
		rType = rType.Elem()
	}

	return rType, wasPtr
}

// DefineVariable enrich the Type by adding field with given name.
// val can be either of the reflect.Type or regular type (i.e. Foo)
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
	if p.planListener != nil {
		p.planListener.OnDefineVariable(name, sType)
	}

	for _, additionalFieldName := range names {
		if err := p.addSelectors("", field, additionalFieldName); err != nil {
			return err
		}
		if p.planListener != nil {
			p.planListener.OnDefineVariable(additionalFieldName, sType)
		}
	}

	return nil
}

func (p *Planner) selector(selector *expr.Select) (*op.Selector, error) {
	resultSelector, call, err := p.matchFirstSelector(selector)
	if err != nil || call == nil {
		return resultSelector, err
	}

	for call != nil {
		resultSelector, call, err = p.matchSelector(call, resultSelector, resultSelector.ID, deref(resultSelector.Type))
		if err != nil {
			return nil, err
		}
	}
	return resultSelector, nil
}

func (p *Planner) matchFirstSelector(selector *expr.Select) (*op.Selector, ast.Expression, error) {
	opSelector, next, err := p.tryMatchFunc(selector.ID, selector.X)
	if err == nil && opSelector != nil {
		return opSelector, next, nil
	}

	resultSelector := p.selectorByName(selector.ID)
	if resultSelector != nil {
		return resultSelector, selector.X, nil
	}

	return nil, nil, nil
}

func (p *Planner) matchSelector(call ast.Expression, resultSelector *op.Selector, selectorId string, parentType reflect.Type) (*op.Selector, ast.Expression, error) {
	switch actual := call.(type) {
	case *expr.Call:
		// Direct method on current selector (e.g., $foo.ToUpper())
		return p.matchFunc(selectorId, actual, resultSelector)

	case *expr.SliceIndex:
		// Indexing into current selector (e.g., $arr[0])
		switch resultSelector.Type.Kind() {
		case reflect.Map:
			mapSelector, err := p.newMapSelector(selectorId, actual, resultSelector)
			if err != nil {
				return nil, nil, err
			}
			return mapSelector, actual.Y, nil
		case reflect.Interface:
			interfaceSelector, err := p.newInterfaceSelector(selectorId, actual, resultSelector)
			if err != nil {
				return nil, nil, err
			}
			return interfaceSelector, actual.Y, nil
		default:
			sliceSelector, err := p.newSliceSelector(selectorId, actual, resultSelector)
			if err != nil {
				return nil, nil, err
			}
			return sliceSelector, actual.Y, nil
		}

	case *expr.Select:
		// If the next segment is a method call on this field (e.g., $Slice.Length()),
		// handle it before resolving the field as a struct field.
		if _, isCall := actual.X.(*expr.Call); isCall {
			if callSelector, callNext, callErr := p.tryMatchCall(actual.X, resultSelector, actual.ID); callSelector != nil || callErr != nil {
				return callSelector, callNext, callErr
			}
		}

		fieldSelector, err := p.selectField(resultSelector, parentType, actual.ID)
		return fieldSelector, actual.X, err
	}

	return resultSelector, nil, nil
}

func deref(rType reflect.Type) reflect.Type {
	seen := map[reflect.Type]bool{}
	for {
		if seen[rType] {
			return rType
		}
		seen[rType] = true
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
	sel.NodeID = p.nextNodeID()
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

	if t == nil {
		return fmt.Errorf("couldn't determine %v type", expr.Name)
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
	aFunc, err := p.Func(prev, methodName, call)
	if err != nil {
		return nil, err
	}

	operands, err := p.selectorOperands(call, prev)
	if err != nil {
		return nil, err
	}

	resultType := aFunc.ResultType
	if resultType == xreflect.InterfaceType {
		actualType, err := p.Functions.TryDetectResultType(prev, methodName, call)
		if err != nil {
			return nil, err
		}
		if actualType != nil {
			resultType = actualType
		}
	}
	accumulator := p.accumulator(resultType)
	newSelector := op.FunctionSelector(selectorId, accumulator.Field, aFunc, prev)
	newSelector.NodeID = p.nextNodeID()
	newSelector.Args = operands
	newSelector.Type = resultType

	if p.planListener != nil {
		var recvType reflect.Type
		if prev != nil {
			recvType = prev.Type
		}
		p.planListener.OnFunctionBind(methodName, aFunc, recvType)
	}

	return newSelector, nil
}

func (p *Planner) Func(prev *op.Selector, methodName string, call *expr.Call) (*op.Func, error) {
	var receiver reflect.Type
	if prev != nil {
		receiver = prev.Type
	}
	aFunc, err := p.Functions.Method(receiver, methodName, call)
	if err == nil {
		return aFunc, nil
	}

	return nil, fmt.Errorf("not found function %v", methodName)
}

func (p *Planner) selectorOperands(call *expr.Call, prev *op.Selector) ([]*op.Operand, error) {
	operands := make([]*op.Operand, 0, len(call.Args)+1)
	if prev != nil {
		operand, err := op.NewExpression(prev).Operand(*p.Control)
		if err != nil {
			return nil, err
		}

		operands = append(operands, operand)
	}

	for i := 0; i < len(call.Args); i++ {
		expression, err := p.compileExpr(call.Args[i])
		if err != nil {
			return nil, err
		}

		operand, err := expression.Operand(*p.Control)
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)
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

func (p *Planner) nextNodeID() int {
	p.nodeSeq++
	return p.nodeSeq
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

func (p *Planner) apply(options []Option) {
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
		case *op.Functions:
			if p.Functions == nil {
				p.Functions = actual
			}
		case Listener:
			p.listener = ParserListener(actual)
		case Adjuster:
			p.adjuster = NodeAdjuster(actual)
		case Policies:
			if reg := (*PolicyRegistry)(actual); reg != nil {
				p.adjuster = reg.AsAdjuster()
			}
		case EvalCfg:
			p.evalConfig = EvaluateConfig(actual)
		case PlanHooks:
			p.planListener = PlannerListener(actual)
		}
	}
}

func (p *Planner) registerConst(i *[]int) {
	p.constants.add(i)
}

func (p *Planner) init(options []Option) {
	p.apply(options)

	if p.Functions == nil {
		ifaces := make([]interface{}, 0, len(options))
		for _, option := range options {
			ifaces = append(ifaces, option)
		}

		p.Functions = op.NewFunctions(ifaces...)
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
		return p.matchFunc(ID, actual, selector)
	case *expr.SliceIndex:
		// If indexing follows a field segment (e.g., $Model.Items[0]), resolve that field first.
		if ID != "" && selector != nil {
			combined := selector.ID + fieldSeparator + ID
			if resolved, ok := p.selectors.ById(combined); ok {
				selector = resolved
				ID = combined
			}
		}

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

	s, err := op.SliceSelector(id, "", sliceOperand, operandExpression, selector)
	if err != nil {
		return nil, err
	}
	s.NodeID = p.nextNodeID()
	return s, nil
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

	s, err := op.NewMapSelector(id, "", mapOperand, keyOperand, selector)
	if err != nil {
		return nil, err
	}
	s.NodeID = p.nextNodeID()
	return s, nil
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

	s, err := op.NewInterfaceSelector(id, "", xOperand, yOperand, selector)
	if err != nil {
		return nil, err
	}
	s.NodeID = p.nextNodeID()
	return s, nil
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

func (p *Planner) tryMatchFunc(ID string, x ast.Expression) (*op.Selector, ast.Expression, error) {
	call, ok := x.(*expr.Call)
	if !ok {
		return nil, nil, nil
	}

	return p.matchFunc(ID, call, nil)
}

func (p *Planner) matchFunc(ID string, actual *expr.Call, selector *op.Selector) (*op.Selector, ast.Expression, error) {
	callSelector, err := p.newFuncSelector(ID, ID, actual, selector)
	if err != nil {
		return nil, nil, err
	}

	return callSelector, actual.X, nil
}
