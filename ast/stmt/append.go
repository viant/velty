package stmt

type Append struct {
	Append string
}

func NewAppend(value string) *Append {
	return &Append{Append: value}
}
