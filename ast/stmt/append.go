package stmt

//Append represents regular text without any template expressions
type Append struct {
	Append string
}

//NewAppend creates new *Append
func NewAppend(value string) *Append {
	return &Append{Append: value}
}
