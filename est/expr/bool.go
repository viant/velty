package expr

import (
	"github.com/viant/velty/est"
	"github.com/viant/velty/internal/ast"
	"unsafe"
)

func computeBinaryBool(token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
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

func computeUnaryBool(token ast.Token, unary *unaryExpr, indirect bool) (est.Compute, error) {
	switch token {
	case ast.NEG:
		if indirect {
			return unary.indirectBoolNeq, nil
		}
		return unary.directBoolNeg, nil
	case "":
		if indirect {
			return unary.indirectBool, nil
		}
		return unary.directBoo, nil
	}

	return nil, errorUnsupported(token, "Bool")
}

func (e *unaryExpr) directBoolNeg(state *est.State) unsafe.Pointer {
	x := e.x.Pointer(state)
	y := est.FalseValuePtr
	if !*(*bool)(x) {
		y = est.TrueValuePtr
	}

	return y
}

func (e *unaryExpr) indirectBoolNeq(state *est.State) unsafe.Pointer {
	x := e.x.Exec(state)
	y := est.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		y = est.TrueValuePtr
	}
	return y
}

func (e *unaryExpr) indirectBool(state *est.State) unsafe.Pointer {
	x := e.x.Exec(state)
	y := est.FalseValuePtr

	if *(*bool)(x) == *(*bool)(y) {
		y = est.TrueValuePtr
	}
	return y
}

func (e *unaryExpr) directBoo(state *est.State) unsafe.Pointer {
	x := e.x.Pointer(state)
	y := est.FalseValuePtr

	if *(*bool)(x) {
		y = est.TrueValuePtr
	}
	return y
}
