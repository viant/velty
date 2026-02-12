package velty

import (
	"github.com/stretchr/testify/require"
	aexpr "github.com/viant/velty/ast/expr"
	"testing"
)

type captureListener struct{ events []Event }

func (c *captureListener) OnEvent(e Event) { c.events = append(c.events, e) }

func TestParserSpansPositions(t *testing.T) {
	tpl := []byte("hello\n$foo.Bar(42)[0]\nworld")
	cap := &captureListener{}
	p := New(Listener(cap))
	_, _, err := p.Compile(tpl)
	require.NoError(t, err)

	// find a Select event with non-zero positions
	var selPos Position
	for _, ev := range cap.events {
		switch ev.Node.(type) {
		case *aexpr.Select:
			if ev.Position.Line != 0 || ev.Position.Col != 0 {
				selPos = ev.Position
			}
		}
	}
	require.NotZero(t, selPos.Line, "selector position should be populated")
}
