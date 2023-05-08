package expr

import (
	"fmt"
	"github.com/viant/velty/ast"
	"github.com/viant/velty/est"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

func computePtr(rType reflect.Type, token ast.Token, binary *binaryExpr, indirect bool) (est.Compute, error) {
	switch rType.Elem().Kind() {
	case reflect.Int, reflect.Int64:
		switch token {
		case ast.GTR:
			return compareIntPtrGreater(binary)
		case ast.GTE:
			return compareIntGreaterOrEqual(binary)
		case ast.LSS:
			return compareIntLess(binary)
		case ast.LEQ:
			return compareIntLessOrEqual(binary)
		case ast.NEQ:
			return compareIntNotEqual(binary)
		case ast.EQ:
			return compareIntPtrEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for int/Int64", token)
		}

	case reflect.Int8:
		switch token {
		case ast.GTR:
			return compareInt8Greater(binary)
		case ast.GTE:
			return compareInt8GreaterOrEqual(binary)
		case ast.LSS:
			return compareInt8Less(binary)
		case ast.LEQ:
			return compareInt8LessOrEqual(binary)
		case ast.EQ:
			return compareInt8Equal(binary)
		case ast.NEQ:
			return compareInt8NotEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for int8", token)
		}
	case reflect.Int16:
		switch token {
		case ast.GTR:
			return compareInt16Greater(binary)
		case ast.GTE:
			return compareInt16GreaterOrEqual(binary)
		case ast.LSS:
			return compareInt16Less(binary)
		case ast.LEQ:
			return compareInt16LessOrEqual(binary)
		case ast.NEQ:
			return compareInt16NotEqual(binary)
		case ast.EQ:
			return compareInt16Equal(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for int16", token)
		}
	case reflect.Uint, reflect.Uint64:
		switch token {
		case ast.GTR:
			return compareUintGreater(binary)
		case ast.GTE:
			return compareUintGreaterOrEqual(binary)
		case ast.LSS:
			return compareUint8Less(binary)
		case ast.LEQ:
			return compareUintLessOrEqual(binary)
		case ast.NEQ:
			return compareUintNotEqual(binary)
		case ast.EQ:
			return compareUintPtrEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for uint/Uint64", token)
		}
	case reflect.Uint8:
		switch token {
		case ast.GTR:
			return compareUInt8Greater(binary)
		case ast.GTE:
			return compareUint8GreaterOrEqual(binary)
		case ast.LSS:
			return compareUint8Less(binary)
		case ast.LEQ:
			return compareUint8LessOrEqual(binary)
		case ast.NEQ:
			return compareUint8NotEqual(binary)
		case ast.EQ:
			return compareUint8Equal(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for uint8", token)
		}
	case reflect.Uint16:
		switch token {
		case ast.GTR:
			return compareUInt16Greater(binary)
		case ast.GTE:
			return compareUint16GreaterOrEqual(binary)
		case ast.LSS:
			return compareUint16Less(binary)
		case ast.LEQ:
			return compareUint16LessOrEqual(binary)
		case ast.NEQ:
			return compareUint16NotEqual(binary)
		case ast.EQ:
			return compareUint16Equal(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for uint16", token)
		}

	case reflect.Int32:
		switch token {
		case ast.GTR:
			return compareInt32Greater(binary)
		case ast.GTE:
			return compareInt32GreaterOrEqual(binary)
		case ast.LSS:
			return compareInt32Less(binary)
		case ast.LEQ:
			return compareInt32LessOrEqual(binary)
		case ast.EQ:
			return compareInt32Equal(binary)
		case ast.NEQ:
			return compareInt32NotEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for int32", token)
		}

	case reflect.Uint32:
		switch token {
		case ast.GTR:
			return compareUInt32Greater(binary)
		case ast.GTE:
			return compareUint32GreaterOrEqual(binary)
		case ast.LSS:
			return compareUint32Less(binary)
		case ast.LEQ:
			return compareUint32LessOrEqual(binary)
		case ast.EQ:
			return compareUint32Equal(binary)
		case ast.NEQ:
			return compareUint32NotEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for int32", token)
		}

	case reflect.Float32:
		switch token {
		case ast.GTR:
			return compareFloat32Greater(binary)
		case ast.GTE:
			return compareFloat32GreaterOrEqual(binary)
		case ast.LSS:
			return compareFloat32PtrLess(binary)
		case ast.LEQ:
			return compareFloat32PtrLessEqual(binary)
		case ast.EQ:
			return compareFloat32PtrEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for float32", token)
		}

	case reflect.Float64:
		switch token {
		case ast.GTR:
			return compareFloat64Greater(binary)
		case ast.GTE:
			return compareFloat64GreaterOrEqual(binary)
		case ast.LSS:
			return compareFloat64PtrLess(binary)
		case ast.LEQ:
			return compareFloat64PtrLessEqual(binary)
		case ast.EQ:
			return compareFloat64PtrEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for float64", token)
		}

	case reflect.Bool:
		switch token {
		case ast.EQ:
			return compareBoolEqual(binary)
		case ast.NEQ:
			return compareBoolNotEqual(binary)
		}
	case reflect.String:
		switch token {
		case ast.EQ:
			return compareStringEqual(binary)
		case ast.NEQ:
			return compareStringNotEqual(binary)
		default:
			return nil, fmt.Errorf("unsupported token %v for string", token)
		}
	}

	return nil, fmt.Errorf("unsupported binary type %v", rType.String())
}

func compareIntLess(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt(xPtr) < xunsafe.AsInt(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareIntLessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt(xPtr) <= xunsafe.AsInt(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareIntGreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt(xPtr) >= xunsafe.AsInt(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareIntPtrGreater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt(xPtr) > xunsafe.AsInt(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt8Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt8(xPtr) > xunsafe.AsInt8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt16Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt16(xPtr) > xunsafe.AsInt16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt32Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt32(xPtr) > xunsafe.AsInt32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUintGreater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint(xPtr) > xunsafe.AsUint(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUInt8Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint8(xPtr) > xunsafe.AsUint8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUInt16Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint16(xPtr) > xunsafe.AsUint16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUInt32Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint32(xPtr) > xunsafe.AsUint32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt8GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt8(xPtr) >= xunsafe.AsInt8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt16GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt16(xPtr) >= xunsafe.AsInt16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt32GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt32(xPtr) >= xunsafe.AsInt32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUintGreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint(xPtr) >= xunsafe.AsUint(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint8GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint8(xPtr) >= xunsafe.AsUint8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint16GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint16(xPtr) >= xunsafe.AsUint16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint32GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint32(xPtr) >= xunsafe.AsUint32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat32Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsFloat32(xPtr) > xunsafe.AsFloat32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat64Greater(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsFloat64(xPtr) > xunsafe.AsFloat64(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat32GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsFloat32(xPtr) >= xunsafe.AsFloat32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat64GreaterOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.TrueValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsFloat64(xPtr) >= xunsafe.AsFloat64(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt8Equal(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt8(xPtr) == xunsafe.AsInt8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt16Equal(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt16(xPtr) == xunsafe.AsInt16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt32Equal(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsInt32(xPtr) == xunsafe.AsInt32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint8Equal(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint8(xPtr) == xunsafe.AsUint8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint16Equal(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint16(xPtr) == xunsafe.AsUint16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint32Equal(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xunsafe.AsUint32(xPtr) == xunsafe.AsUint32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt8Less(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt8(xPtr) < xunsafe.AsInt8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt16Less(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt16(xPtr) < xunsafe.AsInt16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt32Less(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt32(xPtr) < xunsafe.AsInt32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint8Less(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint8(xPtr) < xunsafe.AsUint8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint16Less(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint16(xPtr) < xunsafe.AsUint16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint32Less(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint32(xPtr) < xunsafe.AsUint32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUintLessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint(xPtr) <= xunsafe.AsUint(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt8LessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt8(xPtr) <= xunsafe.AsInt8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint8LessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint8(xPtr) <= xunsafe.AsUint8(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt16LessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt16(xPtr) <= xunsafe.AsInt16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint16LessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint16(xPtr) <= xunsafe.AsUint16(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt32LessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsInt32(xPtr) <= xunsafe.AsInt32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint32LessOrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsUint32(xPtr) <= xunsafe.AsUint32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareStringEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xStr := xunsafe.AsString(xPtr)
		yStr := xunsafe.AsString(yPtr)

		if xStr == yStr {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareStringNotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr == nil || yPtr == nil {
			return est.TrueValuePtr
		}

		xStr := xunsafe.AsString(xPtr)
		yStr := xunsafe.AsString(yPtr)

		if xStr != yStr {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareBoolNotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.TrueValuePtr
		}

		xBool := xunsafe.AsBool(xPtr)
		yBool := xunsafe.AsBool(yPtr)

		if xBool != yBool {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareBoolEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xBool := xunsafe.AsBool(xPtr)
		yBool := xunsafe.AsBool(yPtr)

		if xBool == yBool {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareIntNotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xInt := xunsafe.AsInt(xPtr)
		yInt := xunsafe.AsInt(yPtr)

		if xInt != yInt {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUintNotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xUint := xunsafe.AsUint(xPtr)
		yUint := xunsafe.AsUint(yPtr)

		if xUint != yUint {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint8NotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xUint8 := xunsafe.AsUint8(xPtr)
		yUint8 := xunsafe.AsUint8(yPtr)

		if xUint8 != yUint8 {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt8NotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xInt8 := xunsafe.AsInt8(xPtr)
		yInt8 := xunsafe.AsInt8(yPtr)

		if xInt8 != yInt8 {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt16NotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xInt16 := xunsafe.AsInt16(xPtr)
		yInt16 := xunsafe.AsInt16(yPtr)

		if xInt16 != yInt16 {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint16NotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xUint16 := xunsafe.AsUint16(xPtr)
		yUint16 := xunsafe.AsUint16(yPtr)

		if xUint16 != yUint16 {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareInt32NotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xInt32 := xunsafe.AsInt32(xPtr)
		yInt32 := xunsafe.AsInt32(yPtr)

		if xInt32 != yInt32 {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUint32NotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xUint32 := xunsafe.AsUint32(xPtr)
		yUint32 := xunsafe.AsUint32(yPtr)

		if xUint32 != yUint32 {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat32PtrLess(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsFloat32(xPtr) < xunsafe.AsFloat32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat64PtrLess(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.FalseValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsFloat64(xPtr) < xunsafe.AsFloat64(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat32PtrLessEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsFloat32(xPtr) <= xunsafe.AsFloat32(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat64PtrLessEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if yPtr == nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xPtr != nil && yPtr == nil {
			return est.FalseValuePtr
		}

		if yPtr != nil && xPtr == nil {
			return est.TrueValuePtr
		}

		if xunsafe.AsFloat64(xPtr) <= xunsafe.AsFloat64(yPtr) {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareIntPtrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xVal := xunsafe.AsInt(xPtr)
		yVal := xunsafe.AsInt(yPtr)

		if xVal == yVal {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareUintPtrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xVal := xunsafe.AsUint(xPtr)
		yVal := xunsafe.AsUint(yPtr)

		if xVal == yVal {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat32PtrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xVal := xunsafe.AsFloat32(xPtr)
		yVal := xunsafe.AsFloat32(yPtr)

		if xVal == yVal {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat64PtrEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xVal := xunsafe.AsFloat64(xPtr)
		yVal := xunsafe.AsFloat64(yPtr)

		if xVal == yVal {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat32PtrNotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xVal := xunsafe.AsFloat32(xPtr)
		yVal := xunsafe.AsFloat32(yPtr)

		if xVal != yVal {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}

func compareFloat64PtrNotEqual(binary *binaryExpr) (est.Compute, error) {
	return func(state *est.State) unsafe.Pointer {
		xPtr := binary.x.Exec(state)
		yPtr := binary.y.Exec(state)

		if xPtr == nil || yPtr == nil {
			return est.FalseValuePtr
		}

		xVal := xunsafe.AsFloat64(xPtr)
		yVal := xunsafe.AsFloat64(yPtr)

		if xVal != yVal {
			return est.TrueValuePtr
		}

		return est.FalseValuePtr
	}, nil
}
