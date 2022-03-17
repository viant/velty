package expr

import (
	"fmt"
	"github.com/viant/velty/ast"
)

func errorUnsupported(token ast.Token, dataType string) error {
	return fmt.Errorf("unsupported %v use on %v", token, dataType)
}
