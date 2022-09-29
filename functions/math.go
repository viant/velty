package functions

import (
	"fmt"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"math"
	"reflect"
)

type Math struct{}

func (m Math) Discover(aFunc interface{}) (func(operands []*op.Operand, state *est.State) (interface{}, error), reflect.Type, bool) {
	switch actual := aFunc.(type) {
	case func(_ Math, arg float64) float64:
		return func(operands []*op.Operand, state *est.State) (interface{}, error) {
			if len(operands) != 1 {
				return nil, fmt.Errorf("unexpected number of operands, expected 1, got %v", len(operands))
			}

			return actual(m, *(*float64)(operands[0].Exec(state))), nil
		}, floatType, true
	}

	return nil, nil, false
}

func (m Math) Round(f float64) float64 {
	return math.Round(f)
}

func (m Math) Ceil(f float64) float64 {
	return math.Ceil(f)
}

func (m Math) Abs(f float64) float64 {
	return math.Abs(f)
}

func (m Math) Floor(f float64) float64 {
	return math.Floor(f)
}

func (m Math) Sqrt(f float64) float64 {
	return math.Sqrt(f)
}

func (m Math) Pow(x, y float64) float64 {
	return math.Pow(x, y)
}

func (m Math) Min(x, y float64) float64 {
	return math.Min(x, y)
}

func (m Math) Max(x, y float64) float64 {
	return math.Max(x, y)
}
