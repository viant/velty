package expr

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"unsafe"
)

func computeInt(token ast.Token, binary *binaryExpr) (est.Compute, error) {
	switch token {
	case ast.QUO:
		return binary.intQuo, nil
	case ast.ADD:
		return binary.intAdd, nil
	case ast.SUB:
		return binary.intSub, nil
	case ast.MUL:
		return binary.intMul, nil
	case ast.NEQ:
		return binary.intNeq, nil
	case ast.EQ:
		return binary.intEq, nil
	case ast.GTR:
		return binary.intGtr, nil
	case ast.GTE:
		return binary.intGte, nil
	case ast.LSS:
		return binary.intLss, nil
	case ast.LEQ:
		return binary.intLeq, nil
	}
	return nil, errorUnsupported(token, "Integer")
}

func (b *binaryExpr) intQuo(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) / *(*int)(y)
	return z
}

func (b *binaryExpr) intAdd(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) + *(*int)(y)

	return z
}

func (b *binaryExpr) intSub(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) - *(*int)(y)

	return z
}

func (b *binaryExpr) intMul(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) * *(*int)(y)

	return z
}

func (b *binaryExpr) intEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*int)(x) == *(*int)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) intNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*int)(x) != *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) intGtr(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) > *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) intGte(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) >= *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) intLss(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) < *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) intLeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) <= *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}
