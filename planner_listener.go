package velty

import (
	"github.com/viant/velty/est/op"
	"reflect"
)

// PlannerListener receives hooks during planning/binding (post-parse).
type PlannerListener interface {
	// OnDefineVariable is called when a planner defines a top-level variable.
	OnDefineVariable(name string, rtype reflect.Type)
	// OnSelectorResolved is called when an expression resolves to a selector.
	OnSelectorResolved(sel *op.Selector)
	// OnFunctionBind is called when a function/method is bound.
	OnFunctionBind(name string, f *op.Func, receiverType reflect.Type)
	// OnForEachResolved is called when a foreach item is resolved with types.
	OnForEachResolved(itemName string, itemType, setType reflect.Type)
}
