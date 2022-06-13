package functions

import "strconv"

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
