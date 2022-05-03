package expr

import (
	est "github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeBinaryInt(token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
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

func (b *binaryExpr) indirectIntQuo(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) / *(*int)(y)
	return z
}

func (b *binaryExpr) directIntQuo(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) / *(*int)(y)
	return z
}

func (b *binaryExpr) indirectIntAdd(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) + *(*int)(y)

	return z
}

func (b *binaryExpr) directIntAdd(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) + *(*int)(y)

	return z
}

func (b *binaryExpr) indirectIntSub(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) - *(*int)(y)

	return z
}

func (b *binaryExpr) directIntSub(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) - *(*int)(y)

	return z
}

func (b *binaryExpr) indirectIntMul(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := b.z.Pointer(state)

	*(*int)(z) = *(*int)(x) * *(*int)(y)
	return z
}

func (b *binaryExpr) directIntMul(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := b.z.Pointer(state)
	*(*int)(z) = *(*int)(x) * *(*int)(y)

	return z
}

func (b *binaryExpr) indirectIntEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*int)(x) == *(*int)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) directIntEq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr

	if *(*int)(x) == *(*int)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *binaryExpr) indirectIntNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr

	if *(*int)(x) != *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntNeq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr

	if *(*int)(x) != *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectIntGtr(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) > *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntGtr(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*int)(x) > *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectIntGte(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) >= *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntGte(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*int)(x) >= *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectLss(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) < *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directLss(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*int)(x) < *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) indirectIntLeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := est.FalseValuePtr
	if *(*int)(x) <= *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}

func (b *binaryExpr) directIntLeq(state *est.State) unsafe.Pointer {
	x := b.x.Pointer(state)
	y := b.y.Pointer(state)
	z := est.FalseValuePtr
	if *(*int)(x) <= *(*int)(y) {
		z = est.TrueValuePtr
	}
	return z
}
