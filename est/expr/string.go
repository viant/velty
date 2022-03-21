package expr

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"unsafe"
)

func computeString(token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
	switch token {
	case ast.ADD:
		if indirect {
			return binary.indirectStringAdd, nil
		}
		return binary.directStringAdd, nil
	case ast.EQ:
		if indirect {
			return binary.indirectStringEq, nil
		}
		return binary.directStringEq, nil
	case ast.NEQ:
		if indirect {
			return binary.indirectStringNeq, nil
		}
		return binary.directStringNeq, nil
	}
	return nil, errorUnsupported(token, "string")
}

func (b *binaryExpr) indirectStringAdd(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*string)(z) = *(*string)(x) + *(*string)(y)
	return z
}

func (b *binaryExpr) directStringAdd(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*string)(z) = *(*string)(x) + *(*string)(y)
	return z
}

func (b *binaryExpr) indirectStringEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est.FalseValuePtr
	if *(*string)(x) == *(*string)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directStringEq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)

	z := est.FalseValuePtr
	if *(*string)(x) == *(*string)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectStringNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est.FalseValuePtr
	if *(*string)(x) != *(*string)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directStringNeq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)

	z := est.FalseValuePtr
	if *(*string)(x) != *(*string)(y) {
		z = est.TrueValuePtr
	}

	return z
}
