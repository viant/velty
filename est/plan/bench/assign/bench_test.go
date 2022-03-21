package _if

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/plan"
	"testing"
)

//go:embed template.vm
var template string

var exec *est.Execution
var state *est.State

func init() {
	planner := plan.New(1024)
	var err error
	var benchNewState func() *est.State
	exec, benchNewState, err = planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	state = benchNewState()
}

func Benchmark_Exec_Shared(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		state.Reset()
		exec.Exec(state)
	}
	assert.Equal(b, " 1321321\n abc\n false\n 10000.321", state.Buffer.String())
}
