package est

import "unsafe"

type Compute func(mem *State) unsafe.Pointer

type New func(control Control) Compute
