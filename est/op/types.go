package op

import (
	"github.com/viant/velty/ast/expr"
	"reflect"
)

type MethodResultTyper interface {
	MethodResultType(methodName string, call *expr.Call) (reflect.Type, error)
}
