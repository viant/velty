package velty

import (
	"github.com/viant/velty/est/op"
	"reflect"
)

// CycleDetector tracks types on one registration ancestry, not across siblings.
type CycleDetector struct {
	rType  reflect.Type
	parent *CycleDetector
}

// MultiCycleDetector is retained for callers using Get directly.
type MultiCycleDetector struct{}

func NewCycleDetector(rType reflect.Type) *CycleDetector {
	rType, _ = elemIfNeeded(rType)
	return &CycleDetector{rType: rType}
}

func (c *CycleDetector) Child(rType reflect.Type, _ *op.Selector) (*CycleDetector, bool) {
	rType, _ = elemIfNeeded(rType)
	for ancestor := c; ancestor != nil; ancestor = ancestor.parent {
		if ancestor.rType == rType {
			return ancestor, true
		}
	}
	return &CycleDetector{rType: rType, parent: c}, false
}

func (d *MultiCycleDetector) Get(c *CycleDetector, rType reflect.Type, parentSelector *op.Selector) (*CycleDetector, bool) {
	return c.Child(rType, parentSelector)
}

func (c *CycleDetector) Has(parent *CycleDetector) bool {
	if parent == nil {
		return false
	}
	for ancestor := c; ancestor != nil; ancestor = ancestor.parent {
		if ancestor.rType != nil && ancestor.rType == parent.rType {
			return true
		}
	}
	return false
}
