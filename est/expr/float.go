package expr

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeFloat(token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
	switch token {
	case ast.QUO:
		if indirect {
			return binary.indirectFloatQuo, nil
		}
		return binary.directFloatQuo, nil

	case ast.ADD:
		if indirect {
			return binary.indirectFloatAdd, nil
		}

		return binary.directFloatAdd, nil
	case ast.SUB:
		if indirect {
			return binary.indirectFloatSub, nil
		}

		return binary.directFloatSub, nil
	case ast.MUL:
		if indirect {
			return binary.indirectFloatMul, nil
		}
		return binary.directFloatMul, nil
	case ast.EQ:
		if indirect {
			return binary.indirectFloatEq, nil
		}

		return binary.directFloatEq, nil
	case ast.NEQ:
		if indirect {
			return binary.indirectFloatNeq, nil
		}

		return binary.directFloatNeq, nil
	case ast.GTR:
		if indirect {
			return binary.indirectFloatGtr, nil
		}

		return binary.directFloatGtr, nil
	case ast.GTE:
		if indirect {
			return binary.indirectFloatGte, nil
		}

		return binary.directFloatGte, nil
	case ast.LSS:
		if indirect {
			return binary.indirectFloatLss, nil
		}
		return binary.directFloatLss, nil
	case ast.LEQ:

		return binary.indirectFloatLeq, nil
	}
	return nil, errorUnsupported(token, "Float64")
}

func (b *binaryExpr) indirectFloatQuo(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) / *(*float64)(y)
	return z
}

func (b *binaryExpr) directFloatQuo(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) / *(*float64)(y)
	return z
}

func (b *binaryExpr) indirectFloatAdd(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) + *(*float64)(y)

	return z
}

func (b *binaryExpr) directFloatAdd(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) + *(*float64)(y)

	return z
}

func (b *binaryExpr) indirectFloatSub(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) - *(*float64)(y)

	return z
}

func (b *binaryExpr) directFloatSub(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) - *(*float64)(y)

	return z
}

func (b *binaryExpr) indirectFloatMul(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) * *(*float64)(y)

	return z
}

func (b *binaryExpr) directFloatMul(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) * *(*float64)(y)

	return z
}

func (b *binaryExpr) indirectFloatEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*float64)(x) == *(*float64)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directFloatEq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr

	if *(*float64)(x) == *(*float64)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectFloatNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est.FalseValuePtr
	if *(*float64)(x) != *(*float64)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directFloatNeq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)

	z := est.FalseValuePtr
	if *(*float64)(x) != *(*float64)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectFloatGtr(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) > *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatGtr(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*float64)(x) > *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectFloatGte(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) >= *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatGte(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*float64)(x) >= *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectFloatLss(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) < *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatLss(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*float64)(x) < *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectFloatLeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*float64)(x) <= *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatLeq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*float64)(x) <= *(*float64)(y) {
		z = est.TrueValuePtr
	}
	return z
}
