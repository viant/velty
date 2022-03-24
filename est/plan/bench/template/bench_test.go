package template

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/plan"
	"strings"
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
	initIndirect()
}

func initDirect() {
	vars := map[string]interface{}{
		"Foo": &benchStruct,
	}

	planner := plan.New(1024)
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

	planner := plan.New(1024)
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
	assert.Equal(b, "\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n\n", directState.Buffer.String())
}

func Benchmark_Exec_Indirect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		indirectState.Reset()
		indirectExec.Exec(indirectState)
	}
	assert.Equal(b, "\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n    \n    \n\n    \n    NameLen < 5\n    \nName: FooNameLower: Foo\n    \n    \n\n    \n    NameLen < 5\n    \nName: BarNameLower: Bar\n    \n    \n\n    \n    NameLen < 5\n    \nName: JohnNameLower: John\n    \n    \n\n    \n    NameLen < 5\n    \nName: BobNameLower: Bob\n    \n    \n\n    \n    NameLen < 5\n    \nName: LiamNameLower: Liam\n    \n    \n\n    \n    NameLen < 5\n    \nName: NoahNameLower: Noah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliverNameLower: Oliver\n    \n    \n\n    \n    NameLen > 5\n    \nName: ElijahNameLower: Elijah\n    \n    \n\n    \n    NameLen > 5\n    \nName: OliviaNameLower: Olivia\n    \n    \n\n    \n    NameLen < 5\n    \nName: EmmaNameLower: Emma\n    \n    \n\n    \n    NameLen < 5\n    \nName: AvaNameLower: Ava\n    \n    \n\n    \n    NameLen > 5\n    \nName: CharlotteNameLower: Charlotte\n\n", indirectState.Buffer.String())
}

func Benchmark_NewState_Direct(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		directNewState()
	}
}
