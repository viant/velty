package est

type Execution struct {
	compute Compute
}

func (e *Execution) Exec(stat *State) {
	e.compute(stat)
}

func NewExecution(compute Compute) *Execution {
	return &Execution{compute: compute}
}
