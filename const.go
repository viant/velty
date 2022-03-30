package velty

type constants struct {
	counter   int
	constants map[int]interface{}
}

func (c *constants) add(value interface{}) {
	c.constants[c.counter] = value
	c.counter++
}

func newConstants() *constants {
	return &constants{
		constants: map[int]interface{}{},
	}
}
