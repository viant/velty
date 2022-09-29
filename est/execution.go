package est

import "fmt"

type Execution struct {
	compute Compute
}

func (e *Execution) Exec(stat *State) error {
	e.compute(stat)
	if len(stat.Errors) > 0 {
		return fmt.Errorf("error occured while processing template: %w", stat.Errors[0])
	}

	return nil
}

func NewExecution(compute Compute) *Execution {
	return &Execution{compute: compute}
}
