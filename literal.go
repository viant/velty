package velty

import (
	"fmt"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/internal/ast/expr"
	"reflect"
	"strconv"
	"unsafe"
)

func (p *Planner) literalExpr(literal *expr.Literal) (*op.Expression, error) {
	expr := &op.Expression{}
	switch literal.RType.Kind() {
	case reflect.Int:
		i, _ := strconv.Atoi(literal.Value)
		p.constants.add(&i)
		ptr := unsafe.Pointer(&i)
		expr.Type = reflect.TypeOf(i)
		expr.LiteralPtr = &ptr
	case reflect.Float64:
		f, _ := strconv.ParseFloat(literal.Value, 64)
		p.constants.add(&f)
		expr.Type = reflect.TypeOf(f)
		ptr := unsafe.Pointer(&f)
		expr.LiteralPtr = &ptr
	case reflect.Bool:
		b, _ := strconv.ParseBool(literal.Value)
		p.constants.add(&b)
		expr.Type = reflect.TypeOf(b)
		ptr := unsafe.Pointer(&b)
		expr.LiteralPtr = &ptr
	case reflect.String:
		expr.Type = reflect.TypeOf(literal.Value)
		p.constants.add(&literal.Value)
		ptr := unsafe.Pointer(&literal.Value)
		expr.LiteralPtr = &ptr
	default:
		return nil, fmt.Errorf("invalid literal type: %v", literal.RType.String())
	}
	return expr, nil
}
