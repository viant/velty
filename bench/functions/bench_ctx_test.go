package _if

import (
	"context"
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	est "github.com/viant/velty/est"
	"strings"
	"testing"
)

// A context-aware version of ToUpper
func toUpperWithCtx(ctx context.Context, s string) string {
	return strings.ToUpper(s)
}

//go:embed template_ctx.vm
var templateCtx string

var execCtx *est.Execution
var stateCtx *est.State

func init() {
	var err error
	planner := velty.New(velty.BufferSize(1024))
	if err = planner.DefineVariable("Name", ""); err != nil {
		panic(err)
	}
	if err = planner.RegisterFunction("toUpperWithCtx", toUpperWithCtx); err != nil {
		panic(err)
	}
	if err = planner.RegisterFunction("trim", strings.TrimSpace); err != nil {
		panic(err)
	}
	var benchNewState func() *est.State
	execCtx, benchNewState, err = planner.Compile([]byte(templateCtx))
	if err != nil {
		panic(err)
	}
	stateCtx = benchNewState()
	if err = stateCtx.SetValue("Name", "foo"); err != nil {
		panic(err)
	}
}

func Benchmark_Exec_WithContext(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		stateCtx.Reset()
		_ = execCtx.ExecWithContext(ctx, stateCtx)
	}
	assert.Equal(b, "FOO", strings.TrimSpace(stateCtx.Buffer.String()))
}
