package expr

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"reflect"
	"unsafe"
)

type directBinary struct {
	x *op.Operand
	y *op.Operand
	z *op.Operand
}

func (b *directBinary) quo(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) / *(*int)(y)
	return z
}

func (b *directBinary) add(state *est.State) unsafe.Pointer {
	x := b.x.Exec(state)
	y := b.y.Exec(state)
	z := state.Pointer(*b.z.Offset)
	*(*int)(z) = *(*int)(x) + *(*int)(y)

	fmt.Printf("add %v %v %v\n", *(*int)(z), *(*int)(x), *(*int)(y))
	return z
}

func Binary(token ast.Token, exprs ...*op.Expression) (est.New, error) {

	return func(control est.Control) (est.Compute, error) {
		oprands, err := op.Expressions(exprs).Operands(control)
		if err != nil {
			return nil, err
		}
		fmt.Printf("ADD BINARY: %v\n", token)
		binary := &directBinary{x: oprands[op.X], y: oprands[op.Y], z: oprands[op.Z]}
		switch exprs[0].Type.Kind() {

		case reflect.Int:
			switch token {
			case ast.QUO:
				return binary.quo, nil
			case ast.ADD:
				return binary.add, nil
			}

		case reflect.String:

		}
		return nil, fmt.Errorf("unsupported")
	}, nil
}
