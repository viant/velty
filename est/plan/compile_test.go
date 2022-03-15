package plan_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/plan"
	"testing"
)

func TestPlanner_Compile(t *testing.T) {
	type foo struct {
		Name string
	}

	type address struct {
		Street string
	}

	type employee struct {
		Address address
	}

	var testCases = []struct {
		description string
		template    string
		vars        map[string]interface{}
		expect      string
	}{
		{
			description: "assignment",
			template:    `#set($var1 = 123)$var1`,
			expect:      `123`,
		},
		{
			description: "assign binary expression result",
			template:    `#set($var1 = 12400/100)$var1`,
			expect:      `124`,
		},
		{
			description: "assign binary expression result with select at the right side",
			template:    `#set($var1 = 3 + $num)$var1`,
			expect:      `13`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: "assign binary expression result with select at the left side #1",
			template:    `#set($var1 = $num +  3)$var1`,
			expect:      `13`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: "assign binary expression result with select at the left side #2",
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
			description: `assign binary expression, multiplication`,
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
			expect:      `def`,
		},
		{
			description: "if statement #2",
			template:    `#if(1==1) abc #else def#end`,
			expect:      `abc`,
		},
		{
			description: "if statement #3",
			template:    `#if($var1==$var2) abc #else def#end`,
			expect:      `def`,
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
			expect: `abc`,
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
			expect: `variables are not equalvar1 is bigger than var2`,
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
			expect: `variables are not equalvar2 is bigger than var1`,
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
			expect: `The value of var:0The value of var:1The value of var:2The value of var:3The value of var:4The value of var:5The value of var:6The value of var:7The value of var:8The value of var:9`,
		},
		{
			description: "for statement #1",
			template: `
#foreach($var in $values)
	The value of var: $var
#end
`,
			expect: `The value of var:1The value of var:2The value of var:3The value of var:4The value of var:5The value of var:6The value of var:7The value of var:8The value of var:9`,
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
			template:    `${employee.Address.Street}`,
			expect:      `employee street`,
			vars: map[string]interface{}{
				"employee": employee{Address: address{Street: "employee street"}},
			},
		},
	}
outer:
	//for i, testCase := range testCases[len(testCases)-1:] {
	for i, testCase := range testCases[24:25] {
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
var benchExec *est.Execution
var benchNewState func() *est.State
var benchState *est.State

func init() {
	return
	template := `
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

#for($var3 = 0; $var3 < 110; $var3++) 
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
	benchExec, benchNewState, err = planner.Compile([]byte(template))
	if err != nil {
		fmt.Println(err.Error())
	}

	benchState = benchNewState()
}

func BenchmarkExec(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchExec.Exec(benchState)
		benchState.Reset()
	}
}
