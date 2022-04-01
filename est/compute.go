package est

import "unsafe"

type Compute func(state *State) unsafe.Pointer

type New func(control Control) (Compute, error)

type Computers []New

func (c Computers) New(control Control) ([]Compute, error) {
	var result = make([]Compute, len(c))
	var err error
	for i, n := range c {
		if result[i], err = n(control); err != nil {
			return nil, err
		}
	}
	return result, nil
}
