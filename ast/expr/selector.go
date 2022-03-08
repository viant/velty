package expr

import (
	"fmt"
	"github.com/viant/velty/ast"
	"reflect"
)

type Select struct {
	ID string
}

func (s Select) AddStatement(_ ast.Statement) error {
	return fmt.Errorf("unepxected Add Statement to Select")
}

func (s Select) Type() reflect.Type {
	//TODO: Should selector have type? If no, how do we create i.e. expr.If
	return nil
}
