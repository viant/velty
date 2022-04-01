package expr

import (
	est2 "github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeInt(token ast.Token, binary *binaryExpr, indirect bool) (est2.Compute, error) {
	switch token {
	case ast.QUO:
		if indirect {
			return binary.indirectIntQuo, nil
		}

		return binary.directIntQuo, nil
	case ast.ADD:
		if indirect {
			return binary.indirectIntAdd, nil
		}

		return binary.directIntAdd, nil
	case ast.SUB:
		if indirect {
			return binary.indirectIntSub, nil
		}

		return binary.directIntSub, nil
	case ast.MUL:
		if indirect {
			return binary.indirectIntMul, nil
		}

		return binary.directIntMul, nil
	case ast.NEQ:
		if indirect {
			return binary.indirectIntNeq, nil
		}

		return binary.directIntNeq, nil
	case ast.EQ:
		if indirect {
			return binary.indirectIntEq, nil
		}

		return binary.directIntEq, nil
	case ast.GTR:
		if indirect {
			return binary.indirectIntGtr, nil
		}

		return binary.directIntGtr, nil
	case ast.GTE:
		if indirect {
			return binary.indirectIntGte, nil
		}

		return binary.directIntGte, nil
	case ast.LSS:
		if indirect {
			return binary.indirectLss, nil
		}

		return binary.directLss, nil
	case ast.LEQ:
		if indirect {
			return binary.indirectIntLeq, nil
		}

		return binary.directIntLeq, nil
	}
	return nil, errorUnsupported(token, "Integer")
}

func (b *binaryExpr) indirectIntQuo(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) / *(*int)(y)
	return z
}

func (b *binaryExpr) directIntQuo(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) / *(*int)(y)
	return z
}

func (b *binaryExpr) indirectIntAdd(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) + *(*int)(y)

	return z
}

func (b *binaryExpr) directIntAdd(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) + *(*int)(y)

	return z
}

func (b *binaryExpr) indirectIntSub(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) - *(*int)(y)

	return z
}

func (b *binaryExpr) directIntSub(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) - *(*int)(y)

	return z
}

func (b *binaryExpr) indirectIntMul(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)

	*(*int)(z) = *(*int)(x) * *(*int)(y)
	return z
}

func (b *binaryExpr) directIntMul(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) * *(*int)(y)

	return z
}

func (b *binaryExpr) indirectIntEq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr

	if *(*int)(x) == *(*int)(y) {
		z = est2.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directIntEq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr

	if *(*int)(x) == *(*int)(y) {
		z = est2.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectIntNeq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr

	if *(*int)(x) != *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntNeq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr

	if *(*int)(x) != *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectIntGtr(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*int)(x) > *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntGtr(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*int)(x) > *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectIntGte(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*int)(x) >= *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntGte(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*int)(x) >= *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectLss(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*int)(x) < *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directLss(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*int)(x) < *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectIntLeq(state *est2.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est2.FalseValuePtr
	if *(*int)(x) <= *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntLeq(state *est2.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est2.FalseValuePtr
	if *(*int)(x) <= *(*int)(y) {
		z = est2.TrueValuePtr
	}
	return z
}
