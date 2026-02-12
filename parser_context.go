package velty

import (
	"github.com/viant/velty/ast"
)

// Span represents byte offsets for a node within a source template.
type Span struct {
	Start int
	End   int
}

// Position holds line/column coordinates for a span.
type Position struct {
	Line    int
	Col     int
	EndLine int
	EndCol  int
}

// ExprContextKind classifies where an expression/node appears in the template.
type ExprContextKind uint8

const (
	CtxUnknown ExprContextKind = iota
	// Control flow
	CtxIfCond
	CtxIfBody
	CtxElseIfCond
	CtxElseBody
	CtxForEachCond
	CtxForEachBody
	CtxForLoopInit
	CtxForLoopCond
	CtxForLoopPost
	// Assignment
	CtxSetLHS
	CtxSetRHS
	// Expression vs text
	CtxAppendExpr
	CtxPlainText
	// Function arguments
	CtxFuncArg
	// Evaluate expression
	CtxEvaluate
)

// ExprContext describes the current expression context for a node.
type ExprContext struct {
	Kind   ExprContextKind
	ArgIdx int16 // for CtxFuncArg, -1 otherwise
}

// VarKind classifies a symbol within the current parser context.
type VarKind uint8

const (
	VarUnknown VarKind = iota
	VarLocal
	VarParam
	VarConst
	VarNamespace
	VarFunc
)

// Scope represents a lexical scope with variables and an optional parent scope.
type Scope struct {
	Parent *Scope
	Vars   map[string]interface{}
}

// NewScope creates a new Scope with the given parent.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Vars:   make(map[string]interface{}),
	}
}

// ParserContext holds parser state, including scope stack and node ancestry.
type ParserContext struct {
	// Current variable scope
	Scope *Scope
	// Symbol sets for resolution helpers
	locals     map[string]struct{}
	params     map[string]struct{}
	namespaces map[string]struct{}
	funcs      map[string]struct{}
	consts     map[string]struct{}
	decorators map[string]struct{}
	// occ keeps per-name occurrence counters for selectors/variables.
	occ map[string]int
	// NodeStack tracks ancestor nodes
	NodeStack []ast.Node
	// exprCtxStack tracks expression/statement context
	exprCtxStack []ExprContext
	// Source and position mapping (optional)
	File       string
	Source     []byte
	lineStarts []int
	// Node span registry
	nodeSpans map[ast.Node]Span
	patches   []Patch
	Reporter  *Reporter
	// Evaluate handling configuration
	EvalConfig EvaluateConfig
	EvalDepth  int
}

// CurrentExprContext returns the innermost expression context, if any.
func (p *ParserContext) CurrentExprContext() ExprContext {
	if len(p.exprCtxStack) == 0 {
		return ExprContext{Kind: CtxUnknown, ArgIdx: -1}
	}
	return p.exprCtxStack[len(p.exprCtxStack)-1]
}

// PushExprContext adds a new expression context on the stack.
func (p *ParserContext) PushExprContext(ctx ExprContext) {
	p.exprCtxStack = append(p.exprCtxStack, ctx)
}

// PopExprContext removes the most recent expression context.
func (p *ParserContext) PopExprContext() {
	if len(p.exprCtxStack) == 0 {
		return
	}
	p.exprCtxStack = p.exprCtxStack[:len(p.exprCtxStack)-1]
}

// EnterScope pushes a new nested scope on the context.
func (p *ParserContext) EnterScope() {
	p.Scope = NewScope(p.Scope)
}

// ExitScope pops the current scope, reverting to the parent.
func (p *ParserContext) ExitScope() {
	if p.Scope != nil {
		p.Scope = p.Scope.Parent
	}
}

// SeedSymbols initializes known params, namespaces, functions, and consts.
func (p *ParserContext) SeedSymbols(params, namespaces, funcs, consts []string) {
	if p.params == nil {
		p.params = map[string]struct{}{}
	}
	if p.namespaces == nil {
		p.namespaces = map[string]struct{}{}
	}
	if p.funcs == nil {
		p.funcs = map[string]struct{}{}
	}
	if p.consts == nil {
		p.consts = map[string]struct{}{}
	}
	for _, s := range params {
		p.params[s] = struct{}{}
	}
	for _, s := range namespaces {
		p.namespaces[s] = struct{}{}
	}
	for _, s := range funcs {
		p.funcs[s] = struct{}{}
	}
	for _, s := range consts {
		p.consts[s] = struct{}{}
	}
}

// MarkLocal records a local variable in the current scope and set.
func (p *ParserContext) MarkLocal(name string) {
	if p.locals == nil {
		p.locals = map[string]struct{}{}
	}
	p.locals[name] = struct{}{}
	if p.Scope != nil {
		if p.Scope.Vars == nil {
			p.Scope.Vars = map[string]interface{}{}
		}
		p.Scope.Vars[name] = true
	}
}

// MarkUnexpandRaw marks a name as having UnexpandRaw decorator.
func (p *ParserContext) MarkUnexpandRaw(name string) {
	if p.decorators == nil {
		p.decorators = map[string]struct{}{}
	}
	p.decorators[name] = struct{}{}
}

// IsLocal reports if a symbol is a local variable.
func (p *ParserContext) IsLocal(name string) bool {
	if p.locals != nil {
		if _, ok := p.locals[name]; ok {
			return true
		}
	}
	for s := p.Scope; s != nil; s = s.Parent {
		if s.Vars != nil {
			if _, ok := s.Vars[name]; ok {
				return true
			}
		}
	}
	return false
}

// IsParam reports if a symbol is a planner-defined param (root variable).
func (p *ParserContext) IsParam(name string) bool {
	if p.params == nil {
		return false
	}
	_, ok := p.params[name]
	return ok
}

// IsNamespace reports if a name is a function namespace.
func (p *ParserContext) IsNamespace(name string) bool {
	if p.namespaces == nil {
		return false
	}
	_, ok := p.namespaces[name]
	return ok
}

// IsStandaloneFunc reports if a name is a registered standalone function.
func (p *ParserContext) IsStandaloneFunc(name string) bool {
	if p.funcs == nil {
		return false
	}
	_, ok := p.funcs[name]
	return ok
}

// IsConst reports if a symbol is a known constant.
func (p *ParserContext) IsConst(name string) bool {
	if p.consts == nil {
		return false
	}
	_, ok := p.consts[name]
	return ok
}

// HasUnexpandRaw reports if a name is marked with UnexpandRaw decorator.
func (p *ParserContext) HasUnexpandRaw(name string) bool {
	if p.decorators == nil {
		return false
	}
	_, ok := p.decorators[name]
	return ok
}

// VarKind reports the classification of a symbol within this context.
func (p *ParserContext) VarKind(name string) VarKind {
	if name == "" {
		return VarUnknown
	}
	if p.IsLocal(name) {
		return VarLocal
	}
	if p.IsParam(name) {
		return VarParam
	}
	if p.IsConst(name) {
		return VarConst
	}
	if p.IsNamespace(name) {
		return VarNamespace
	}
	if p.IsStandaloneFunc(name) {
		return VarFunc
	}
	return VarUnknown
}

// BumpOccurrence increments and returns the 1-based occurrence index for name.
func (p *ParserContext) BumpOccurrence(name string) int {
	if name == "" {
		return 0
	}
	if p.occ == nil {
		p.occ = make(map[string]int)
	}
	next := p.occ[name] + 1
	p.occ[name] = next
	return next
}

// PushNode records a node entry into the context stack.
func (p *ParserContext) PushNode(node ast.Node) {
	p.NodeStack = append(p.NodeStack, node)
}

// PopNode removes the most recent node from the context stack.
func (p *ParserContext) PopNode() {
	if len(p.NodeStack) > 0 {
		p.NodeStack = p.NodeStack[:len(p.NodeStack)-1]
	}
}

// CurrentNode returns the most recently pushed node or nil if empty.
func (p *ParserContext) CurrentNode() ast.Node {
	if len(p.NodeStack) == 0 {
		return nil
	}
	return p.NodeStack[len(p.NodeStack)-1]
}

// InitSource initializes source bytes and precomputes line starts.
func (p *ParserContext) InitSource(file string, src []byte) {
	p.File = file
	p.Source = src
	p.lineStarts = make([]int, 0, 128)
	p.lineStarts = append(p.lineStarts, 0)
	for i, b := range src {
		if b == '\n' {
			p.lineStarts = append(p.lineStarts, i+1)
		}
	}
	if p.nodeSpans == nil {
		p.nodeSpans = make(map[ast.Node]Span)
	}
}

// SetSpan registers a span for a node.
func (p *ParserContext) SetSpan(n ast.Node, s Span) {
	if p.nodeSpans == nil {
		p.nodeSpans = make(map[ast.Node]Span)
	}
	p.nodeSpans[n] = s
}

// GetSpan retrieves a span for a node, if any.
func (p *ParserContext) GetSpan(n ast.Node) (Span, bool) {
	if p.nodeSpans == nil {
		return Span{}, false
	}
	s, ok := p.nodeSpans[n]
	return s, ok
}

// ResolvePosition converts a span to line/column positions.
func (p *ParserContext) ResolvePosition(s Span) Position {
	pos := Position{}
	if p.lineStarts == nil || len(p.lineStarts) == 0 || s.Start < 0 || s.End < 0 {
		return pos
	}
	// find start line
	pos.Line, pos.Col = p.byteToLineCol(s.Start)
	pos.EndLine, pos.EndCol = p.byteToLineCol(s.End)
	return pos
}

func (p *ParserContext) byteToLineCol(b int) (int, int) {
	// binary search over lineStarts
	lo, hi := 0, len(p.lineStarts)-1
	for lo <= hi {
		mid := (lo + hi) >> 1
		if p.lineStarts[mid] <= b {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	lineIdx := lo - 1
	if lineIdx < 0 {
		lineIdx = 0
	}
	lineStart := p.lineStarts[lineIdx]
	return lineIdx + 1, b - lineStart + 1
}

// AddPatches appends patches to the context accumulator.
func (p *ParserContext) AddPatches(ps ...Patch) {
	p.patches = append(p.patches, ps...)
}

// Patches returns accumulated patches.
func (p *ParserContext) Patches() []Patch { return p.patches }

// EvaluateMode controls how #evaluate is handled during traversal.
type EvaluateMode int

const (
	EvalOpaque EvaluateMode = iota
	EvalInspect
	EvalRewrite
)

// EvaluateSafety defines safety constraints for handling #evaluate content.
type EvaluateSafety struct {
	StringLiteralsOnly bool
	MaxBytes           int // 0 = unlimited
}

// EvaluateConfig defines behavior for #evaluate traversal.
type EvaluateConfig struct {
	Mode       EvaluateMode
	DepthLimit int // 0 = unlimited
	Safety     EvaluateSafety
	Whitelist  []string // allowed selector IDs for evaluated content
	Blacklist  []string // disallowed selector IDs
}
