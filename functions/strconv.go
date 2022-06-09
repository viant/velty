package functions

import "strconv"

type Strconv struct{}

func (s Strconv) Itoa(i int) string {
	return strconv.Itoa(i)
}

func (s Strconv) Atoi(val string) (int, error) {
	return strconv.Atoi(val)
}
