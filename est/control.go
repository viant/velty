package est

const (
	// ControlHasBreak marks templates that include break statements.
	ControlHasBreak Control = 1 << iota
	// ControlInLoop marks a scope that is inside #for/#foreach.
	ControlInLoop
)

// Control represents execution control flags like uses continue, uses break
type Control uint8

func (c Control) HasBreak() bool {
	return c&ControlHasBreak != 0
}

func (c Control) InLoop() bool {
	return c&ControlInLoop != 0
}
