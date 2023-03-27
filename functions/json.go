package functions

import (
	"encoding/json"
	"fmt"
	"github.com/viant/velty/ast/expr"
	"io"
	"reflect"
)

const unmarshalIntoMethodName = "UnmarshalInto"

type JSON struct {
	parser TypeParser
}

func (n *JSON) MethodResultType(methodName string, call *expr.Call) (reflect.Type, error) {
	if methodName == unmarshalIntoMethodName {
		args := call.Args
		if len(args) != 2 {
			return nil, fmt.Errorf("unexpected number of args, expected 2 got %v", len(args))
		}

		asLiteral, ok := args[1].(*expr.Literal)
		if !ok {
			return nil, nil
		}

		return n.parser(asLiteral.Value)
	}

	return nil, nil
}

func NewJSON(typeParser TypeParser) *JSON {
	return &JSON{parser: typeParser}
}

func (n *JSON) Marshal(any interface{}) (string, error) {
	marshal, err := json.Marshal(any)
	return string(marshal), err
}

func (n *JSON) UnmarshalInto(data interface{}, typeRepresentation string) (interface{}, error) {
	if n.parser == nil {
		return nil, fmt.Errorf("can't unmarshall into %v due to not specified TypeParser", typeRepresentation)
	}

	bytes, err := asBytes(data)
	if err != nil {
		return nil, err
	}

	rType, err := n.parser(typeRepresentation)
	if err != nil {
		return nil, err
	}

	rValue := reflect.New(rType)
	if err = json.Unmarshal(bytes, rValue.Interface()); err != nil {
		return nil, err
	}

	return rValue.Elem().Interface(), nil

}

func asBytes(data interface{}) ([]byte, error) {
	switch actual := data.(type) {
	case string:
		return []byte(actual), nil
	case *string:
		if actual != nil {
			return []byte(*actual), nil
		}
		return []byte(""), nil
	case []byte:
		return actual, nil
	case io.Reader:
		return io.ReadAll(actual)
	}

	return nil, fmt.Errorf("couldn't convert %T into []byte", data)
}
