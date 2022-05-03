package expr

import (
	est "github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeBinaryString(token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
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
	z := b.z.Pointer(state)
	*(*string)(z) = *(*string)(x) + *(*string)(y)
	return z
}

func (b *binaryExpr) directStringAdd(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := b.z.Pointer(state)
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
