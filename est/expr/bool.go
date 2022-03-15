package expr

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"unsafe"
)

func computeBool(token ast.Token, binary *directBinary) (est.Compute, error) {
	switch token {
	case ast.EQ:
		return binary.boolEq, nil
	case ast.NEQ:
		return binary.boolNeq, nil
	}

	return nil, errorUnsupported(token, "Bool")
}

func (b *directBinary) boolEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *directBinary) boolNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*bool)(x) != *(*bool)(y) {
		z = est.TrueValuePtr
	}
	return z
}