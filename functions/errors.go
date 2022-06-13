package functions

import "fmt"

type Errors struct {
}

func (e Errors) RegisterError(message string) (string, error) {
	return "", fmt.Errorf(message)
}
