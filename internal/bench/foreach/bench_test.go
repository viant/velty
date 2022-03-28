package _if

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	"github.com/viant/velty/internal/est"
	"testing"
)

//go:embed template.vm
var template string

var directExec *est.Execution
var directState *est.State
var directNewState func() *est.State

var indirectExec *est.Execution
var indirectState *est.State
var benchStruct = Foo{
	Values: Values{
		Ints: []int{1000, 2500, 43245, 2145532, 12321, 543124214325, 23241321, 534214, 3251, 343531423, 54, 432},
	},
}

type Values struct {
	Ints []int
	ID   int
}

type Foo struct {
	Values Values
	ID     int
}

func init() {
	initDirect()
	initIndirect()
}

func initDirect() {
	vars := map[string]interface{}{
		"Foo": &benchStruct,
	}

	planner := velty.New(velty.BufferSize(1024))

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	directExec, directNewState, err = planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	directState = directNewState()
	for key, value := range vars {
		if err := directState.SetValue(key, value); err != nil {
			fmt.Println(err)
		}
	}
}

func initIndirect() {
	type Foo struct {
		Values *Values
		ID     int
	}

	vars := map[string]interface{}{
		"Foo": &Foo{
			Values: &benchStruct.Values,
		},
	}

	planner := velty.New(velty.BufferSize(1024))

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	var benchNewState func() *est.State
	indirectExec, benchNewState, err = planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	indirectState = benchNewState()
	for key, value := range vars {
		if err := indirectState.SetValue(key, value); err != nil {
			fmt.Println(err)
		}
	}
}

func Benchmark_Exec_Direct(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		directState.Reset()
		directExec.Exec(directState)
	}
	assert.Equal(b, "\n    1000\n\n    2500\n\n    43245\n\n    2145532\n\n    12321\n\n    543124214325\n\n    23241321\n\n    534214\n\n    3251\n\n    343531423\n\n    54\n\n    432\n\n", directState.Buffer.String())
}

func Benchmark_Exec_Indirect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		indirectState.Reset()
		indirectExec.Exec(indirectState)
	}
	assert.Equal(b, "\n    1000\n\n    2500\n\n    43245\n\n    2145532\n\n    12321\n\n    543124214325\n\n    23241321\n\n    534214\n\n    3251\n\n    343531423\n\n    54\n\n    432\n\n", indirectState.Buffer.String())
}

func Benchmark_NewState_Direct(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		directNewState()

	}
}
