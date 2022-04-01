package expr

import (
	est2 "github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeBool(token ast.Token, binary *binaryExpr, indirect bool) (est2.Compute, error) {
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

func (b *binaryExpr) indirectBoolEq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directBoolEq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectBoolNeq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*bool)(x) != *(*bool)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directBoolNeq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*bool)(x) != *(*bool)(y) {
		z = est2.TrueValuePtr
	}
	return z
}
