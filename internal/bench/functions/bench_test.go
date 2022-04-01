package _if

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	est2 "github.com/viant/velty/est"
	"strings"
	"testing"
)

//go:embed template.vm
var template string

var exec *est2.Execution
var state *est2.State

func init() {
	var err error
	planner := velty.New(velty.BufferSize(1024))
	if err = planner.DefineVariable("Name", ""); err != nil {
		panic(err)
	}

	if err = planner.RegisterFunction("toUpperCase", strings.ToUpper); err != nil {
		panic(err)
	}

	if err = planner.RegisterFunction("trim", strings.TrimSpace); err != nil {
		panic(err)
	}

	var benchNewState func() *est2.State
	exec, benchNewState, err = planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	state = benchNewState()
	if err = state.SetValue("Name", "   foo    "); err != nil {
		panic(err)
	}
}

func Benchmark_Exec_Shared(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		state.Reset()
		exec.Exec(state)
	}
	assert.Equal(b, "FOO", state.Buffer.String())
}
