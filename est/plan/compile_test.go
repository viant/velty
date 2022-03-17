package plan_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/plan"
	"strconv"
	"testing"
)

func TestPlanner_Compile(t *testing.T) {
	type foo struct {
		Name string
	}

	type Values struct {
		IntValue     int
		StringValue  string
		BooleanValue bool
		FloatValue   float64
		Index        int
		Values       []string
	}

	type employee struct {
		Address Values
	}

	type Department struct {
		Address *Values
	}

	values := &Values{
		StringValue:  "employee street",
		IntValue:     123456789,
		BooleanValue: true,
		FloatValue:   125.43,
		Index:        10,
		Values:       []string{"Var1", "Var2", "Var3", "Var4"},
	}
	department := &Department{
		Address: values,
	}

	var testCases = []struct {
		description string
		template    string
		vars        map[string]interface{}
		expect      string
	}{
		{
			description: "assign int",
			template:    `#set($var1 = 123)$var1`,
			expect:      `123`,
		},
		{
			description: "assign binary expression result",
			template:    `#set($var1 = 12400/100)$var1`,
			expect:      `124`,
		},
		{
			description: "assign binary expression result with selector at the right side",
			template:    `#set($var1 = 3 + $num)$var1`,
			expect:      `13`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: "assign binary expression result with selector at the left side #1",
			template:    `#set($var1 = $num +  3)$var1`,
			expect:      `13`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: "assign binary expression result with selector at the left side #2",
			template:    `#set($var1 = $num - 3)$var1`,
			expect:      `7`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: `assign binary expression with precedence evaluation #1`,
			template:    `#set( $var1 = 2 * 2 + 3 )$var1`,
			expect:      "7",
		},
		{
			description: `assign binary expression with precedence evaluation #2`,
			template:    `#set( $var1 = 2 * 2 * 2 + 3 * 2 )$var1`,
			expect:      "14",
		},
		{
			description: `assign binary expression with precedence evaluation #3`,
			template:    `#set( $var1 = 2 * (2 + 8) * 2 + 3 * 2 )$var1`,
			expect:      "46",
		},
		{
			description: `assign binary expression, multiplication #4`,
			template:    `#set( $var1 = 2.5 * 2.5)$var1`,
			expect:      "6.25",
		},
		{
			description: `assign binary expression, concatenation #5`,
			template:    `#set( $var1 = "abc" + "cdef")$var1`,
			expect:      "abccdef",
		},
		{
			description: `assign binary expression, comparison #6`,
			template:    `#set( $var1 = "abc" == "cdef")$var1`,
			expect:      "false",
		},
		{
			description: `assign binary expression, comparison #7`,
			template:    `#set( $var1 = "abc" == "abc")$var1`,
			expect:      "true",
		},
		{
			description: `assign binary expression, comparison #8`,
			template:    `#set( $var1 = 1.5 !=  1.5)$var1`,
			expect:      "false",
		},
		{
			description: `assign binary expression, comparison #9`,
			template:    `#set( $var1 = 1.5 ==  1.5)$var1`,
			expect:      "true",
		},
		{
			description: `assign binary expression, comparison #10`,
			template:    `#set( $var1 = 1 ==  1)$var1`,
			expect:      "true",
		},
		{
			description: `assign binary expression, both side selector #10`,
			template:    `#set( $var1 = $num1 +  $num2)$var1`,
			expect:      "25",
			vars: map[string]interface{}{
				"num1": 10,
				"num2": 15,
			},
		},
		{
			description: "if statement #1",
			template:    `#if(1==2) abc #else def#end`,
			expect:      ` def`,
		},
		{
			description: "if statement #2",
			template:    `#if(1==1) abc #else def#end`,
			expect:      ` abc `,
		},
		{
			description: "if statement #3",
			template:    `#if($var1==$var2) abc #else def#end`,
			expect:      ` def`,
			vars: map[string]interface{}{
				"var1": 1,
				"var2": 2,
			},
		},
		{
			description: "if statement #4",
			template: `
#if($var1==$var2) 
	abc 
#else 
	def
#end`,
			expect: "\n \n\tabc \n",
			vars: map[string]interface{}{
				"var1": 1,
				"var2": 1,
			},
		},
		{
			description: "if statement #5",
			template: `
#if($var1 =!= $var2)
	variables are not equal
	#if($var1 > $var2)
		var1 is bigger than var2
	#elseif($var2 > $var1)
		var2 is bigger than var1
	#else
		never happen
	#end
#end
`,
			expect: "\n\n\tvariables are not equal\n\t\n\t\tvar1 is bigger than var2\n\t\n\n",
			vars: map[string]interface{}{
				"var1": 10,
				"var2": 5,
			},
		},
		{
			description: "if statement #5",
			template: `
#if($var1 =! $var2)
	variables are not equal
	#if($var1 > $var2)
		var1 is bigger than var2
	#elseif($var2 > $var1)
		var2 is bigger than var1
	#else
		never happen
	#end
#end
`,
			expect: "\n\n\tvariables are not equal\n\t\n\t\tvar2 is bigger than var1\n\t\n\n",
			vars: map[string]interface{}{
				"var1": 1,
				"var2": 5,
			},
		},
		{
			description: "for statement #1",
			template: `
#for($var = 1; $var < 10; $var++)
	The value of var: $var
#end
`,
			expect: "\n\n\tThe value of var: 1\n\n\tThe value of var: 2\n\n\tThe value of var: 3\n\n\tThe value of var: 4\n\n\tThe value of var: 5\n\n\tThe value of var: 6\n\n\tThe value of var: 7\n\n\tThe value of var: 8\n\n\tThe value of var: 9\n\n",
		},
		{
			description: "for statement #2",
			template: `
#foreach($var in $values)
	The value of var: $var
#end
`,
			expect: "\n\n\tThe value of var: 1\n\n\tThe value of var: 2\n\n\tThe value of var: 3\n\n\tThe value of var: 4\n\n\tThe value of var: 5\n\n\tThe value of var: 6\n\n\tThe value of var: 7\n\n\tThe value of var: 8\n\n\tThe value of var: 9\n\n",
			vars: map[string]interface{}{
				"values": []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		{
			description: "objects  #1",
			template:    `${foo.Name}`,
			expect:      `Foo name`,
			vars: map[string]interface{}{
				"foo": foo{Name: "Foo name"},
			},
		},
		{
			description: "objects  #2",
			template:    `${employee.Address.StringValue}`,
			expect:      `employee street`,
			vars: map[string]interface{}{
				"employee": employee{Address: *values},
			},
		},
		{
			description: "objects  #3",
			template:    `${employee.Address.StringValue}`,
			expect:      `employee street`,
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #4",
			template:    `${employee.Address.IntValue}`,
			expect:      `123456789`,
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #5",
			template:    `${employee.Address.BooleanValue}`,
			expect:      `true`,
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #6",
			template:    `${employee.Address.FloatValue}`,
			expect:      `125.43`,
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #7",
			template: `
#set($abc = ${employee.Address.StringValue})
$abc
`,
			expect: "\n\nemployee street\n",
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #8",
			template: `
#set($abc = ${employee.Address.FloatValue})
$abc
`,
			expect: "\n\n125.43\n",
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #9",
			template: `
#set($abc = ${employee.Address.FloatValue} + ${employee.Address.FloatValue})
$abc
`,
			expect: "\n\n250.86\n",
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #10",
			template: `
#foreach($value in $employee.Address.Values )
	$value ;
#end
`,
			expect: "\n\n\tVar1 ;\n\n\tVar2 ;\n\n\tVar3 ;\n\n\tVar4 ;\n\n",
			vars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #11",
			template: `
#for($var = 1; $var <= 4; $var++ )
	Var${var} ;
#end
`,
			expect: "\n\n\tVar1 ;\n\n\tVar2 ;\n\n\tVar3 ;\n\n\tVar4 ;\n\n",
			vars: map[string]interface{}{
				"employee": department,
			},
		},
	}
outer:
	//for i, testCase := range testCases[len(testCases)-1:] {
	for i, testCase := range testCases {
		fmt.Printf("Running testcase: %v\n", i)
		planner := plan.New(8192)

		for k, v := range testCase.vars {
			err := planner.DefineVariable(k, v)
			if !assert.Nil(t, err, testCase.description) {
				continue outer
			}
		}
		exec, newState, err := planner.Compile([]byte(testCase.template))
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		state := newState()
		for k, v := range testCase.vars {
			err = state.SetValue(k, v)
			if !assert.Nil(t, err, testCase.description+" var "+k) {
				continue outer
			}
		}
		exec.Exec(state)
		output := state.Buffer.Bytes()
		assert.Equal(t, testCase.expect, string(output), testCase.description)
	}

}

// Benchmarks
type benchData struct {
	execution  *est.Execution
	newState   func() *est.State
	benchState *est.State
	template   string
	variables  map[string]interface{}
}

var directBenchData *benchData
var indirectBenchData *benchData

func init() {
	initDirectBench()
	initIndirectBench()
}

func initDirectBench() {
	template := `
#if($var1 =! $var2)
	variables are not equal
#if($var1 > $var2)
		var1 is bigger than var2
	#elseif($var2 > $var1)
		var2 is bigger than var1
	#else
		never happen
	#end
#end

#for($var3 = 0; $var3 < 100; $var3++) 
		varValue: abc
#end
`
	vars := map[string]interface{}{
		"var1": 10,
		"var2": 5,
	}

	planner := plan.New(8192)

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	benchExec, benchNewState, err := planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	directBenchData = &benchData{
		execution:  benchExec,
		newState:   benchNewState,
		benchState: benchNewState(),
	}
}

func initIndirectBench() {
	indirectTemplate := `
#if($foo.Values.Var1 != $foo.Values.Var2)
	variables are not equal
#if($foo.Values.Var1 > $foo.Values.Var2)
		var1 is bigger than var2
	#elseif($foo.Values.Var2 > $foo.Values.Var1)
		var2 is bigger than var1
	#else
		never happen
	#end
#end

#foreach($var3 in $foo.Values.Data) 
		variable: $var3
#end
`
	type Values struct {
		Var1 int
		Var2 int

		Data []string
	}

	type Foo struct {
		Values *Values
		id     int
	}

	values := make([]string, 100)
	for i := 0; i < len(values); i++ {
		values[i] = "var" + strconv.Itoa(i+1)
	}
	foo := &Foo{
		Values: &Values{
			Var1: 10,
			Var2: 5,
			Data: values,
		},
	}

	vars := map[string]interface{}{
		"foo": foo,
	}

	benchFoo = foo

	planner := plan.New(8192)

	for k, v := range vars {
		err := planner.DefineVariable(k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	var err error
	benchExec, benchNewState, err := planner.Compile([]byte(indirectTemplate))
	if err != nil {
		fmt.Println(err.Error())
	}

	state := benchNewState()
	_ = state.SetValue("foo", foo)
	indirectBenchData = &benchData{
		execution:  benchExec,
		newState:   benchNewState,
		benchState: state,
		template:   indirectTemplate,
		variables:  vars,
	}
}

var benchFoo interface{}

func BenchmarkExec_Direct(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		directBenchData.execution.Exec(directBenchData.benchState)
		directBenchData.benchState.Reset()
	}
}

func BenchmarkExec_Indirect(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		indirectBenchData.execution.Exec(indirectBenchData.benchState)
		indirectBenchData.benchState.Reset()
	}
}
