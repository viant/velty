package expr

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"unsafe"
)

func computeBool(token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
	switch token {
	case ast.EQ:
		if indirect {
			return binary.indirectBoolEq, nil
		}
		return binary.directBoolEq, nil
	case ast.NEQ:
		if indirect {
			return binary.indirectBoolNeq, nil
		}
		return binary.directBoolNeq, nil

	}

	return nil, errorUnsupported(token, "Bool")
}

func (b *binaryExpr) indirectBoolEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directBoolEq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectBoolNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*bool)(x) != *(*bool)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directBoolNeq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*bool)(x) != *(*bool)(y) {
		z = est.TrueValuePtr
	}
	return z
}
