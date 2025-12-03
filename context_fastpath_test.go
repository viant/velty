package velty

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// ctx + []byte -> string
func b2s(ctx context.Context, b []byte) string {
	s := strings.ToUpper(string(b))
	if v := ctx.Value("sfx"); v != nil {
		s += v.(string)
	}
	return s
}

// map variants
func mcount(ctx context.Context, m map[string]string) int { return len(m) }
func mhasA(ctx context.Context, m map[string]string) bool { v, ok := m["a"]; return ok && v != "" }
func mjoin(ctx context.Context, m map[string]string) string {
	// join keys in sorted order for determinism
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if len(keys) == 2 && keys[0] > keys[1] {
		keys[0], keys[1] = keys[1], keys[0]
	}
	return strings.Join(keys, ",")
}

func Test_ContextFastPath_BytesToString(t *testing.T) {
	p := New()
	_ = p.DefineVariable("Bin", []byte{})
	_ = p.RegisterFunction("B2S", b2s)

	tpl := `$B2S($Bin)`
	exec, newState, err := p.Compile([]byte(tpl))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	_ = s.SetValue("Bin", []byte("go"))
	s.SetContext(context.WithValue(context.Background(), "sfx", "!"))
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "GO!", s.Buffer.String())
}

func Test_ContextFastPath_MapVariants(t *testing.T) {
	p := New()
	_ = p.DefineVariable("M", map[string]string{})
	_ = p.RegisterFunction("MCount", mcount)
	_ = p.RegisterFunction("MHasA", mhasA)
	_ = p.RegisterFunction("MJoin", mjoin)

	tpl := `$MCount($M)|$MHasA($M)|$MJoin($M)`
	exec, newState, err := p.Compile([]byte(tpl))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	_ = s.SetValue("M", map[string]string{"a": "x", "b": "y"})
	assert.Nil(t, exec.Exec(s))
	// MCount=2, MHasA=true, MJoin="a,b"
	assert.Equal(t, "2|true|a,b", s.Buffer.String())
}
