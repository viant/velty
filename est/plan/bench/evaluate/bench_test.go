package _if

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/plan"
	"testing"
)

//go:embed parent_template.vm
var parentTemplate string

//go:embed template.vm
var template string

var directExec *est.Execution
var directState *est.State

var indirectExec *est.Execution
var indirectState *est.State
var benchStruct = Foo{
	Values: Values{
		Ints: Ints{
			Values: make([]int, 10),
			Size:   10,
		},
	},
}

type Values struct {
	Ints Ints
	ID   int
}

type Ints struct {
	ID     int
	Values []int
	Size   int
}

type Foo struct {
	Values Values
	ID     int
}

func init() {
	for i := 0; i < len(benchStruct.Values.Ints.Values); i++ {
		benchStruct.Values.Ints.Values[i] = i
	}

	initDirect()
	initIndirect()
}

func initDirect() {
	vars := map[string]interface{}{
		"Template": template,
		"Foo":      &benchStruct,
	}

	planner := plan.New(1024, 1)

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	var benchNewState func() *est.State
	directExec, benchNewState, err = planner.Compile([]byte(parentTemplate))
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
		"Template": template,
		"Foo": &Foo{
			Values: &benchStruct.Values,
		},
	}

	planner := plan.New(1024, 1)

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	var benchNewState func() *est.State
	indirectExec, benchNewState, err = planner.Compile([]byte(parentTemplate))
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
	assert.Equal(b, "\n\nSize > abc\n\n\n\n0\n\n1\n\n2\n\n3\n\n4\n\n5\n\n6\n\n7\n\n8\n\n9\n\n", directState.Buffer.String())
}

func Benchmark_Exec_Indirect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		indirectState.Reset()
		indirectExec.Exec(indirectState)
	}
	assert.Equal(b, "\n\nSize > abc\n\n\n\n0\n\n1\n\n2\n\n3\n\n4\n\n5\n\n6\n\n7\n\n8\n\n9\n\n", indirectState.Buffer.String())
}
