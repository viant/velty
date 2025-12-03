package velty

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// plain context-aware function
func withCtxToUpper(ctx context.Context, s string) string {
	// include a suffix from ctx for visibility
	if v := ctx.Value("sfx"); v != nil {
		return stringsToUpper(s) + v.(string)
	}
	return stringsToUpper(s)
}

// helper to avoid importing strings in test to keep style consistent
func stringsToUpper(s string) string {
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] = b[i] - 32
		}
	}
	return string(b)
}

type ctxReceiver struct{}

func (ctxReceiver) Prefix(ctx context.Context, s string) string {
	p, _ := ctx.Value("pfx").(string)
	return p + s
}

func Test_Context_With_Function(t *testing.T) {
	p := New()
	_ = p.RegisterFunction("ToUpperWithCtx", withCtxToUpper)
	exec, newState, err := p.Compile([]byte(`$ToUpperWithCtx("go")`))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	s.SetContext(context.WithValue(context.Background(), "sfx", "!"))
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "GO!", s.Buffer.String())
}

func Test_Context_With_Method(t *testing.T) {
	p := New()
	_ = p.DefineVariable("R", ctxReceiver{})
	exec, newState, err := p.Compile([]byte(`$R.Prefix("bar")`))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	s.SetContext(context.WithValue(context.Background(), "pfx", "foo:"))
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "foo:bar", s.Buffer.String())
}

func Test_Context_ExecWithContext(t *testing.T) {
	p := New()
	_ = p.RegisterFunction("ToUpperWithCtx", withCtxToUpper)
	exec, newState, err := p.Compile([]byte(`$ToUpperWithCtx("go")`))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	ctx := context.WithValue(context.Background(), "sfx", "?")
	assert.Nil(t, exec.ExecWithContext(ctx, s))
	assert.Equal(t, "GO?", s.Buffer.String())
}
