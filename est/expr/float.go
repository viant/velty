package expr

import (
	est2 "github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeFloat(token ast.Token, binary *binaryExpr, indirect bool) (est2.Compute, error) {
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

func (b *binaryExpr) indirectFloatQuo(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) / *(*float64)(y)
	return z
}

func (b *binaryExpr) directFloatQuo(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) / *(*float64)(y)
	return z
}

func (b *binaryExpr) indirectFloatAdd(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) + *(*float64)(y)

	return z
}

func (b *binaryExpr) directFloatAdd(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) + *(*float64)(y)

	return z
}

func (b *binaryExpr) indirectFloatSub(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) - *(*float64)(y)

	return z
}

func (b *binaryExpr) directFloatSub(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) - *(*float64)(y)

	return z
}

func (b *binaryExpr) indirectFloatMul(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) * *(*float64)(y)

	return z
}

func (b *binaryExpr) directFloatMul(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*float64)(z) = *(*float64)(x) * *(*float64)(y)

	return z
}

func (b *binaryExpr) indirectFloatEq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr

	if *(*float64)(x) == *(*float64)(y) {
		z = est2.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directFloatEq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr

	if *(*float64)(x) == *(*float64)(y) {
		z = est2.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectFloatNeq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est2.FalseValuePtr
	if *(*float64)(x) != *(*float64)(y) {
		z = est2.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directFloatNeq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)

	z := est2.FalseValuePtr
	if *(*float64)(x) != *(*float64)(y) {
		z = est2.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectFloatGtr(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) > *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatGtr(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) > *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectFloatGte(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) >= *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatGte(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) >= *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectFloatLss(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) < *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatLss(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) < *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectFloatLeq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) <= *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directFloatLeq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*float64)(x) <= *(*float64)(y) {
		z = est2.TrueValuePtr
	}
	return z
}
