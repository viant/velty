package template

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	est "github.com/viant/velty/est"
	"strings"
	"testing"
)

//go:embed template.vm
var template string

//go:embed template_no_functions.vm
var templateWithoutFunc string

var directExec *est.Execution
var directState *est.State
var directNewState func() *est.State

var noFuncExec *est.Execution
var noFuncState *est.State
var noFuncNewState func() *est.State

var benchStruct = Foo{
	Values: Values{
		Names: generateStrings(),
	},
}

func generateStrings() []string {
	values := []string{"Foo", "Bar", "John", "Bob", "Liam", "Noah", "Oliver", "Elijah", "Olivia", "Emma", "Ava", "Charlotte"}
	result := make([]string, 0)
	for i := 0; i < 10; i++ {
		result = append(result, values...)
	}

	return result
}

type Values struct {
	Names []string
	ID    int
}

type Foo struct {
	Values Values
	ID     int
}

func init() {
	initDirect()
	initTemplateWithoutFunc()
}

func initDirect() {
	vars := map[string]interface{}{
		"Foo": &benchStruct,
	}

	planner := velty.New(velty.BufferSize(10000))
	if err := planner.RegisterFunction("len", func(s string) int {
		return len(s)
	}); err != nil {
		panic(err)
	}

	if err := planner.RegisterFunction("trim", strings.TrimSpace); err != nil {
		panic(err)
	}

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

func initTemplateWithoutFunc() {
	vars := map[string]interface{}{
		"Foo": &benchStruct,
	}

	planner := velty.New(velty.BufferSize(10000))
	if err := planner.RegisterFunction("len", func(s string) int {
		return len(s)
	}); err != nil {
		panic(err)
	}

	if err := planner.RegisterFunction("trim", strings.TrimSpace); err != nil {
		panic(err)
	}

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	noFuncExec, noFuncNewState, err = planner.Compile([]byte(templateWithoutFunc))
	if err != nil {
		fmt.Println(err.Error())
	}

	noFuncState = directNewState()
	for key, value := range vars {
		if err := noFuncState.SetValue(key, value); err != nil {
			fmt.Println(err)
		}
	}
}

func Benchmark_Exec_Velty(b *testing.B) {
	var ns *est.State
	var err error
	b.ReportAllocs()
	ns = directNewState()
	for i := 0; i < b.N; i++ {
		ns.Reset()
		err = ns.SetValue("Foo", &benchStruct)
		directExec.Exec(ns)
	}
	assert.Nil(b, err)
	assert.Equal(b, "\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Foo\nNameLower: Foo\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bar\nNameLower: Bar\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: John\nNameLower: John\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Bob\nNameLower: Bob\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Liam\nNameLower: Liam\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Noah\nNameLower: Noah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Oliver\nNameLower: Oliver\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Elijah\nNameLower: Elijah\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Olivia\nNameLower: Olivia\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Emma\nNameLower: Emma\n\n    \n    \n\n    \n    NameLen < 5\n    \nName: Ava\nNameLower: Ava\n\n    \n    \n\n    \n    NameLen > 5\n    \nName: Charlotte\nNameLower: Charlotte\n\n\n", ns.Buffer.String())
}

func Benchmark_ExecFuncLess_Velty(b *testing.B) {
	var ns *est.State
	var err error
	b.ReportAllocs()

	ns = noFuncNewState()
	for i := 0; i < b.N; i++ {
		ns.Reset()
		err = ns.SetValue("Foo", &benchStruct)
		noFuncExec.Exec(ns)
	}
	assert.Nil(b, err)
	assert.Equal(b, "\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n    \n\n    \n    Name = Foo\n    \nName: Foo\nNameAssigned: Foo\n\n    \n\n    \n    Name = Bar\n    \nName: Bar\nNameAssigned: Bar\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: John\nNameAssigned: John\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Bob\nNameAssigned: Bob\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Liam\nNameAssigned: Liam\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Noah\nNameAssigned: Noah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Oliver\nNameAssigned: Oliver\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Elijah\nNameAssigned: Elijah\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Olivia\nNameAssigned: Olivia\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Emma\nNameAssigned: Emma\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Ava\nNameAssigned: Ava\n\n    \n\n    \n    Name != Foo && Name != Bar\n    \nName: Charlotte\nNameAssigned: Charlotte\n\n\n", ns.Buffer.String())
}

func Benchmark_NewState_Direct(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		directNewState()
	}
}
