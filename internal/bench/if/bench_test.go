package _if

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	est2 "github.com/viant/velty/est"
	"testing"
)

//go:embed template.vm
var template string

var directExec *est2.Execution
var directState *est2.State

var indirectExec *est2.Execution
var indirectState *est2.State
var benchStruct = Foo{
	Values: Values{
		Ints: Ints{
			Var1: 100000,
			Var2: 5000,
		},
	},
}

type Ints struct {
	Var1 int
	Var2 int
}

type Values struct {
	Ints Ints
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
	var benchNewState func() *est2.State
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
	type Values struct {
		Ints *Ints
		ID   int
	}

	type Foo struct {
		Values *Values
		ID     int
	}

	vars := map[string]interface{}{
		"Foo": &Foo{
			Values: &Values{
				Ints: &benchStruct.Values.Ints,
			},
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
	var benchNewState func() *est2.State
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
	assert.Equal(b, "\nvar1 is not equal var2\n    \n    and var1 is greater than var2\n    \n", directState.Buffer.String())
}

func Benchmark_Exec_Indirect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		indirectState.Reset()
		indirectExec.Exec(indirectState)
	}
	assert.Equal(b, "\nvar1 is not equal var2\n    \n    and var1 is greater than var2\n    \n", indirectState.Buffer.String())
}
