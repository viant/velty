package op

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"github.com/viant/xunsafe/converter"
	"reflect"
	"unsafe"
)

type Slice struct {
	XSlice       *xunsafe.Slice
	IndexOperand *Operand
	SliceOperand *Operand
	ToInter      converter.UnifyFn
}

func (s *Slice) Exec(slicePtr unsafe.Pointer, state *est.State) unsafe.Pointer {
	sliceLen := s.XSlice.Len(slicePtr)

	indexPtr := s.IndexOperand.Exec(state)
	if indexPtr == nil {
		return nil
	}

	index := *(*int)(indexPtr)
	if sliceLen <= index {
		panic(fmt.Sprintf("index out of range [%v] with length %v", index, sliceLen))
	}

	elemPtr := s.XSlice.PointerAt(slicePtr, uintptr(index))
	// If the slice element is a pointer type (e.g., []*T), dereference the slot
	// and return the underlying pointer so field access works on the element.
	if s.SliceOperand != nil && s.SliceOperand.Type != nil {
		if s.SliceOperand.Type.Elem().Kind() == reflect.Ptr {
			return *(*unsafe.Pointer)(elemPtr)
		}
	}
	return elemPtr

}
