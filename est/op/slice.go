package op

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"unsafe"
)

type Slice struct {
	XSlice       *xunsafe.Slice
	IndexOperand *Operand
	SliceOperand *Operand
	ToInter      func(pointer unsafe.Pointer) int
}

func (s *Slice) Exec(state *est.State) unsafe.Pointer {
	slicePtr := s.SliceOperand.Exec(state)
	sliceLen := s.XSlice.Len(slicePtr)

	indexPtr := s.IndexOperand.Exec(state)
	index := s.ToInter(indexPtr)

	if sliceLen <= index {
		panic(fmt.Sprintf("index out of range [%v] with length %v", index, sliceLen))
	}

	return s.XSlice.PointerAt(slicePtr, uintptr(index))

}
