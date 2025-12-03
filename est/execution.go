package est

import (
	"context"
	"fmt"
)

type Execution struct {
	compute      Compute
	PanicOnError bool
}

func (e *Execution) Exec(stat *State) (err error) {
	if e.PanicOnError {
		defer func() {
			panicErr := recover()
			if panicErr != nil {
				asTemplateErr, ok := panicErr.(TemplateError)
				if !ok {
					panic(panicErr)
				}
				err = asTemplateErr
			}
		}()
	}

	e.compute(stat)
	if len(stat.Errors) > 0 {
		return fmt.Errorf("error occured while processing template: %w", stat.Errors[0])
	}

	return err
}

func NewExecution(compute Compute) *Execution {
	return &Execution{compute: compute}
}

// ExecWithContext executes with the provided context set on the state.
// Backward compatible: underlying Exec is invoked after setting state context.
func (e *Execution) ExecWithContext(ctx context.Context, stat *State) error {
	stat.Ctx = ctx
	return e.Exec(stat)
}
