package velty

import (
	"github.com/viant/velty/functions"
)

//Option represents Planner generic option
type Option interface{}

//BufferSize represents initial size of the buffer
type BufferSize int

//CacheSize represents cache size in case of the dynamic template evaluation
type CacheSize int

//EscapeHTML escapes HTML in passed variables.
type EscapeHTML bool

//PanicOnError panics and recover when first error returned.
type PanicOnError bool

//TypeParser parses type string representation into reflect.Type
type TypeParser = functions.TypeParser
