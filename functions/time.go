package functions

import (
	"reflect"
	"time"
)

var (
	Now      = time.Now
	timeType = reflect.TypeOf(time.Time{})
)

type Time struct {
}

func (t Time) DiscoverInterfaces(aFunc interface{}) (func(args ...interface{}) (interface{}, error), reflect.Type, bool) {
	switch actual := aFunc.(type) {
	case func(_ Time) time.Time:
		return func(operands ...interface{}) (interface{}, error) {
			aTime := actual(t)

			return aTime, nil
		}, timeType, true

	case func() time.Time:
		return func(operands ...interface{}) (interface{}, error) {
			aTime := actual()

			return aTime, nil
		}, timeType, true
	}

	return nil, nil, false
}

func (t Time) Now() time.Time {
	return Now()
}
