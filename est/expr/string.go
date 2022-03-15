package expr

import (
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"unsafe"
)

func computeString(token ast.Token, binary *directBinary) (est.Compute, error) {
	switch token {
	case ast.ADD:
		return binary.stringAdd, nil
	case ast.EQ:
		return binary.stringEq, nil
	case ast.NEQ:
		return binary.stringNeq, nil
	}
	return nil, errorUnsupported(token, "string")
}

func (b *directBinary) stringAdd(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*string)(z) = *(*string)(x) + *(*string)(y)
	return z
}

func (b *directBinary) stringEq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est.FalseValuePtr
	if *(*string)(x) == *(*string)(y) {
		z = est.TrueValuePtr
	}

	return z
}

func (b *directBinary) stringNeq(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)

	z := est.FalseValuePtr
	if *(*string)(x) != *(*string)(y) {
		z = est.TrueValuePtr
	}

	return z
}
