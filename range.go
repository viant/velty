package velty

import (
	"github.com/viant/velty/est/op"
	aexpr "github.com/viant/velty/internal/ast/expr"
	"reflect"
	"strconv"
	"unsafe"
)

func (p *Planner) compileRange(actual *aexpr.Range) (*op.Expression, error) {
	begin, err := strconv.Atoi(actual.X.Value)
	if err != nil {
		return nil, err
	}

	finish, err := strconv.Atoi(actual.Y.Value)
	if err != nil {
		return nil, err
	}

	var aSlice []int
	if begin < finish {
		aSlice = make([]int, finish-begin)
		for i := 0; i < len(aSlice); i++ {
			aSlice[i] = begin + i
		}
	} else {
		aSlice = make([]int, begin-finish)
		for i := 0; i < len(aSlice); i++ {
			aSlice[i] = begin - i
		}
	}

	p.registerConst(&aSlice)
	slicePtr := unsafe.Pointer(&aSlice)
	return &op.Expression{
		LiteralPtr: &slicePtr,
		Type:       reflect.SliceOf(actual.Type()),
	}, nil
}
