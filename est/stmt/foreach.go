package stmt

import (
	est2 "github.com/viant/velty/est"
	op2 "github.com/viant/velty/est/op"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type ForEach struct {
	Block est2.Compute

	Item *op2.Operand
	X    *op2.Operand

	*xunsafe.Slice
}

func (e *ForEach) compute(state *est2.State) unsafe.Pointer {
	xPtr := state.Pointer(e.X.Sel.Offset)
	l := e.Slice.Len(xPtr)

	var resultPtr unsafe.Pointer
	for i := 0; i < l; i++ {
		v := e.Slice.ValueAt(xPtr, i)
		e.Item.Sel.Set(state.MemPtr, v)
		resultPtr = e.Block(state)
	}

	return resultPtr
}

func (e *ForEach) computePtr(state *est2.State) unsafe.Pointer {
	xPtr := state.Pointer(e.X.Sel.Offset)
	l := e.Slice.Len(xPtr)

	var resultPtr unsafe.Pointer
	for i := 0; i < l; i++ {
		v := e.Slice.ValuePointerAt(xPtr, i)
		e.Item.Sel.SetValue(state.MemPtr, v)
		resultPtr = e.Block(state)
	}

	return resultPtr
}

func (e *ForEach) computeIndirectPtr(state *est2.State) unsafe.Pointer {
	xPtr := e.X.Exec(state)
	l := e.Slice.Len(xPtr)

	var resultPtr unsafe.Pointer
	for i := 0; i < l; i++ {
		v := e.Slice.ValuePointerAt(xPtr, i)
		e.Item.Sel.SetValue(state.MemPtr, v)
		resultPtr = e.Block(state)
	}

	return resultPtr
}

func (e *ForEach) computeIndirect(state *est2.State) unsafe.Pointer {
	xPtr := e.X.Exec(state)
	l := e.Slice.Len(xPtr)

	var resultPtr unsafe.Pointer
	for i := 0; i < l; i++ {
		v := e.Slice.ValueAt(xPtr, i)
		e.Item.Sel.SetValue(state.MemPtr, v)
		resultPtr = e.Block(state)
	}

	return resultPtr
}

func (e *ForEach) computeLiteral(state *est2.State) unsafe.Pointer {
	xPtr := *e.X.LiteralPtr
	l := e.Slice.Len(xPtr)

	var resultPtr unsafe.Pointer
	for i := 0; i < l; i++ {
		v := e.Slice.ValueAt(xPtr, i)
		e.Item.Sel.SetValue(state.MemPtr, v)
		resultPtr = e.Block(state)
	}

	return resultPtr
}

func ForEachLoop(block est2.New, itemExpr *op2.Expression, sliceExpr *op2.Expression) (est2.New, error) {
	return func(control est2.Control) (est2.Compute, error) {
		aSlice, err := sliceExpr.Operand(control)
		if err != nil {
			return nil, err
		}

		loop := &ForEach{}
		loop.Block, err = block(control)
		if err != nil {
			return nil, err
		}

		loop.Slice = xunsafe.NewSlice(aSlice.Type)
		loop.X = aSlice

		loop.Item, err = itemExpr.Operand(control)
		if err != nil {
			return nil, err
		}

		switch loop.Slice.Elem().Kind() {
		case reflect.Ptr:
			if loop.X.Sel != nil && loop.X.Sel.Indirect {
				return loop.computeIndirectPtr, nil
			}
			return loop.computePtr, nil
		default:
			if loop.X.Sel != nil && loop.X.Sel.Indirect {
				return loop.computeIndirect, nil
			}

			if loop.X.Sel != nil {
				return loop.compute, nil
			}
			return loop.computeLiteral, nil

		}
	}, nil
}
