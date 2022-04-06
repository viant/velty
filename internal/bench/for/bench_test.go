package _if

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	est "github.com/viant/velty/est"
	"testing"
)

//go:embed template.vm
var template string

var directExec *est.Execution
var directState *est.State

var indirectExec *est.Execution
var indirectState *est.State
var benchStruct = Foo{
	Values: Values{
		Ints: Ints{
			Size: 10,
		},
	},
}

type Values struct {
	Ints Ints
	ID   int
}

type Ints struct {
	ID   int
	Size int
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
	var benchNewState func() *est.State
	directExec, benchNewState, err = planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	directState = benchNewState()
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
	assert.Equal(b, "\n    0\n\n    1\n\n    2\n\n    3\n\n    4\n\n    5\n\n    6\n\n    7\n\n    8\n\n    9\n\n", directState.Buffer.String())
}

func Benchmark_Exec_Indirect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		indirectState.Reset()
		indirectExec.Exec(indirectState)
	}
	assert.Equal(b, "\n    0\n\n    1\n\n    2\n\n    3\n\n    4\n\n    5\n\n    6\n\n    7\n\n    8\n\n    9\n\n", indirectState.Buffer.String())
}
