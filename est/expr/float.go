package expr

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"unsafe"
)

func computeFloat(token ast.Token, binary *binaryExpr) (est.Compute, error) {
	switch token {
	case ast.QUO:
		return binary.floatQuo, nil
	case ast.ADD:
		return binary.floatAdd, nil
	case ast.SUB:
		return binary.floatSub, nil
	case ast.MUL:
		return binary.floatMul, nil
	case ast.EQ:
		return binary.floatEq, nil
	case ast.NEQ:
		return binary.floatNeq, nil
	case ast.GTR:
		return binary.floatGtr, nil
	case ast.GTE:
		return binary.floatGte, nil
	case ast.LSS:
		return binary.floatLss, nil
	case ast.LEQ:
		return binary.floatLeq, nil
	}
	return nil, errorUnsupported(token, "Float64")
}

func (b *binaryExpr) floatQuo(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) / *(*float64)(y)
	return z
}

func (b *binaryExpr) floatAdd(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) + *(*float64)(y)

	return z
}

func (b *binaryExpr) floatSub(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) - *(*float64)(y)

	return z
}

func (b *binaryExpr) floatMul(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) * *(*float64)(y)

	return z
}

func (b *binaryExpr) floatEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*float64)(x) == *(*float64)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) floatNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est.FalseValuePtr
	if *(*float64)(x) != *(*float64)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) floatGtr(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) > *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) floatGte(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) >= *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) floatLss(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) < *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) floatLeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) <= *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}
