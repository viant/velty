package velty

import (
	"github.com/viant/velty/est/op"
	"reflect"
)

type (
	CycleDetector struct {
		rType          reflect.Type
		allTypes       map[reflect.Type]*MultiCycleDetector
		parent         *CycleDetector
		asString       string
		parentSelector *op.Selector
	}

	MultiCycleDetector struct {
		index    map[reflect.Type]int
		children []*CycleDetector
	}
)

func NewCycleDetector(rType reflect.Type) *CycleDetector {
	return newCycleDetector(nil, nil, rType, nil)
}

func newCycleDetector(parent *CycleDetector, state map[reflect.Type]*MultiCycleDetector, rType reflect.Type, parentSelector *op.Selector) *CycleDetector {
	if state == nil {
		state = map[reflect.Type]*MultiCycleDetector{}
	}

	rType, _ = elemIfNeeded(rType)
	return &CycleDetector{
		allTypes:       state,
		parent:         parent,
		rType:          rType,
		asString:       rType.String(),
		parentSelector: parentSelector,
	}
}

func (c *CycleDetector) Child(rType reflect.Type, parentSelector *op.Selector) (next *CycleDetector, hasCycle bool) {
	cycle := c.getChildrenCycleHolder(rType)
	return cycle.Get(c, rType, parentSelector)
}

func (d *MultiCycleDetector) Get(c *CycleDetector, rType reflect.Type, parentSelector *op.Selector) (*CycleDetector, bool) {
	detector, parent := d.getOrCreate(c, rType, parentSelector)

	return detector, parent != nil
}

func (d *MultiCycleDetector) getOrCreate(c *CycleDetector, rType reflect.Type, parentSelector *op.Selector) (current *CycleDetector, parent *CycleDetector) {
	i, ok := d.index[rType]
	if ok {
		next := d.children[i]
		if next.Has(next) {
			return c, next
		}

		return next, nil
	}

	detector := newCycleDetector(c, c.allTypes, rType, parentSelector)
	d.index[rType] = len(d.children)
	d.children = append(d.children, detector)

	return detector, nil
}

func (c *CycleDetector) getChildrenCycleHolder(rType reflect.Type) *MultiCycleDetector {
	cycle, ok := c.allTypes[rType]
	if ok {
		return cycle
	}

	cycle = &MultiCycleDetector{
		index: map[reflect.Type]int{},
	}

	c.allTypes[rType] = cycle
	return cycle
}

func (c *CycleDetector) Has(parent *CycleDetector) bool {
	curr := c.parent
	for curr != nil {
		if curr.rType == parent.rType && curr.rType != nil {
			return true
		}

		curr = curr.parent
	}

	return false
}
