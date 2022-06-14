package functions

import (
	"fmt"
	"strconv"
)

type Strconv struct{}

func (s Strconv) Itoa(i int) string {
	return strconv.Itoa(i)
}

func (s Strconv) Atoi(val string) (int, error) {
	return strconv.Atoi(val)
}

func (s Strconv) ParseFloat(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func (s Strconv) ParseBool(val string) (bool, error) {
	return strconv.ParseBool(val)
}

func (s Strconv) ParseUint(val string) (uint64, error) {
	return strconv.ParseUint(val, 10, 64)
}

func (s Strconv) AsFloat(value interface{}) (float64, error) {
	switch actual := value.(type) {
	case float64:
		return actual, nil
	case string:
		return strconv.ParseFloat(actual, 64)
	case int:
		return float64(actual), nil
	}
	return 0, fmt.Errorf("unconvertable value %v to float64", value)
}
