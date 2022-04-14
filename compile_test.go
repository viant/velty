package velty_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	"github.com/viant/velty/est"
	"strings"
	"testing"
)

type bar struct {
	Name string
}

func (b *bar) UpperCase() string {
	return strings.ToUpper(b.Name)
}

func (b *bar) Concat(values ...string) string {
	return strings.Join(append([]string{b.Name}, values...), " ")
}

func TestPlanner_Compile(t *testing.T) {
	type Foo struct {
		Name string
		ID   int
	}

	type Values struct {
		IntValue     int `velty:"name=Int"`
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
		ID      int
	}

	type ValuesHolder struct {
		Values `velty:"prefix=VARIABLES_"`
	}

	type FooWrapper struct {
		Foo
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

	type Tagged struct {
		ID   int `velty:"-"`
		Name string
	}

	taggedStruct := &Tagged{ID: 100}
	var testCases = []testdata{
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
			definedVars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: "assign binary expression result with selector at the left side #1",
			template:    `#set($var1 = $num +  3)$var1`,
			expect:      `13`,
			definedVars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			description: "assign binary expression result with selector at the left side #2",
			template:    `#set($var1 = $num - 3)$var1`,
			expect:      `7`,
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
				"var1": 10,
				"var2": 5,
			},
		},
		{
			description: "if statement #6",
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
				"values": []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
		},
		{
			description: "objects  #1",
			template:    `${foo.Name}`,
			expect:      `Foo name`,
			definedVars: map[string]interface{}{
				"foo": Foo{Name: "Foo name"},
			},
		},
		{
			description: "objects  #2",
			template:    `${employee.Address.StringValue}`,
			expect:      `employee street`,
			definedVars: map[string]interface{}{
				"employee": employee{Address: *values},
			},
		},
		{
			description: "objects  #3",
			template:    `${employee.Address.StringValue}`,
			expect:      `employee street`,
			definedVars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #4",
			template:    `${employee.Address.Int}`,
			expect:      `123456789`,
			definedVars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #5",
			template:    `${employee.Address.BooleanValue}`,
			expect:      `true`,
			definedVars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "objects  #6",
			template:    `${employee.Address.FloatValue}`,
			expect:      `125.43`,
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
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
			definedVars: map[string]interface{}{
				"employee": department,
			},
		},
		{
			description: "evaluate template in runtime",
			template:    `#evaluate(${foo_template})`,
			expect:      `Var1: 1000, Var2: 13213`,
			definedVars: map[string]interface{}{
				"foo_template": `Var1: $var1, Var2: $var2`,
				"var1":         1000,
				"var2":         13213,
			},
		},
		{
			description: "nil #1",
			expect:      "",
			template:    `${Var.Address.StringValue}`,
			definedVars: map[string]interface{}{
				"Var": &Department{},
			},
		},
		{
			description: "nil #2",
			expect:      "0",
			template:    `${Var.Address.Int}`,
			definedVars: map[string]interface{}{
				"Var": &Department{},
			},
		},
		{
			description: "tags #1",
			expect:      "123456789",
			template:    `${VARIABLES_Int}`,
			embeddedVars: map[string]interface{}{
				"ValuesHolder": ValuesHolder{Values{IntValue: 123456789}},
			},
		},
		{
			description: "func #2",
			expect:      "FOO",
			template:    `${Name.toUpperCase()}`,
			definedVars: map[string]interface{}{
				"Name": "foo",
			},
			functions: map[string]interface{}{
				"toUpperCase": strings.ToUpper,
			},
		},
		{
			description: "func #3",
			expect:      "FOO",
			template:    `${Name.toUpperCase().trim()}`,
			definedVars: map[string]interface{}{
				"Name": "     foo        ",
			},
			functions: map[string]interface{}{
				"toUpperCase": strings.ToUpper,
				"trim":        strings.TrimSpace,
			},
		},
		{
			description: "tags #1",
			template:    `$Tagged.ID`,
			expectError: true,
			definedVars: map[string]interface{}{
				"Tagged": taggedStruct,
			},
		},
		{
			description: "range over slice of structs #1",
			template: `
#foreach ($foo in $foos) 
	$foo.Name 
#end`,
			definedVars: map[string]interface{}{
				"foos": []*Foo{
					{
						Name: "Foo name",
					},
					{
						Name: "Another name",
					},
				},
			},
			expect: "\n \n\tFoo name \n \n\tAnother name \n",
		},
		{
			description: "range over slice of structs #2",
			template: `
#foreach ($foo in $foos) 
	$foo.Name 
#end`,
			definedVars: map[string]interface{}{
				"foos": []Foo{
					{
						Name: "Foo name",
					},
					{
						Name: "Another name",
					},
				},
			},
			expect: "\n \n\tFoo name \n \n\tAnother name \n",
		},
		{
			description: "escape HTML",
			template:    `$FOO`,
			definedVars: map[string]interface{}{
				`FOO`: `<script>alert()</script>`,
			},
			options: []velty.Option{velty.EscapeHTML(true)},
			expect:  "&lt;script&gt;alert()&lt;/script&gt;",
		},
		{
			description: "asc range",
			template: `
#foreach($var in [-10...10]) 
	$var 
#end
`,
			expect: "\n \n\t-10 \n \n\t-9 \n \n\t-8 \n \n\t-7 \n \n\t-6 \n \n\t-5 \n \n\t-4 \n \n\t-3 \n \n\t-2 \n \n\t-1 \n \n\t0 \n \n\t1 \n \n\t2 \n \n\t3 \n \n\t4 \n \n\t5 \n \n\t6 \n \n\t7 \n \n\t8 \n \n\t9 \n\n",
		},
		{
			description: "dsc range",
			template: `
#foreach($var in [10...-10]) 
	$var 
#end
`,
			expect: "\n \n\t10 \n \n\t9 \n \n\t8 \n \n\t7 \n \n\t6 \n \n\t5 \n \n\t4 \n \n\t3 \n \n\t2 \n \n\t1 \n \n\t0 \n \n\t-1 \n \n\t-2 \n \n\t-3 \n \n\t-4 \n \n\t-5 \n \n\t-6 \n \n\t-7 \n \n\t-8 \n \n\t-9 \n\n",
		},
		{
			description: "method receiver",
			template:    `$bar.UpperCase()`,
			definedVars: map[string]interface{}{
				"bar": &bar{
					Name: "bar",
				},
			},
			expect: "BAR",
		},
		{
			description: "method receiver with function calls",
			template:    `$bar.Concat($foo, $var.toUpperCase(), "abcdef")`,
			definedVars: map[string]interface{}{
				"bar": &bar{
					Name: "bar",
				},
				"foo": "fooName",
				"var": "value",
			},
			functions: map[string]interface{}{
				"toUpperCase": strings.ToUpper,
			},
			expect: "bar fooName VALUE abcdef",
		},
		{
			description: "evaluate with non-pointer embed",
			template:    `#evaluate($template)`,
			definedVars: map[string]interface{}{
				"template": "$Name",
			},
			embeddedVars: map[string]interface{}{
				"FooWrapper": FooWrapper{
					Foo{
						Name: "abc",
					},
				},
			},
			expect: "abc",
		},
		{
			description: "unary neg",
			template:    `#if(!$boolValue) abc #else def #end`,
			definedVars: map[string]interface{}{
				"boolValue": true,
			},
			expect: ` def `,
		},
		{
			description: "unary",
			template:    `#if($boolValue) abc #else def #end`,
			definedVars: map[string]interface{}{
				"boolValue": true,
			},
			expect: ` abc `,
		},
		{
			description: "selector as placeholder",
			template:    `$foo`,
			expect:      `$foo`,
		},
		{
			description: "selector block as placeholder",
			template:    `${foo}`,
			expect:      `${foo}`,
		},
		{
			description: "binary &&",
			template:    `#if((1==1) && (2==2)) abc #else def #end`,
			expect:      ` abc `,
		},
		{
			description: "binary ||",
			template:    `#if((1==1) || (2==2)) abc #else def #end`,
			expect:      ` abc `,
		},
	}

	//for i, testCase := range testCases[len(testCases)-1:] {
	for i, testCase := range testCases {
		fmt.Printf("Running testcase: %v\n", i)
		exec, state, err := testCase.init(t)
		if testCase.expectError {
			assert.NotNil(t, err, testCase.description)
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		exec.Exec(state)
		output := state.Buffer.Bytes()
		assert.Equal(t, testCase.expect, string(output), testCase.description)
	}
}

type testdata struct {
	description  string
	template     string
	definedVars  map[string]interface{}
	embeddedVars map[string]interface{}
	functions    map[string]interface{}
	expectError  bool
	expect       string
	options      []velty.Option
}

func (d *testdata) init(t *testing.T) (*est.Execution, *est.State, error) {
	options := []velty.Option{velty.BufferSize(8192)}
	if len(d.options) > 0 {
		options = append(options, d.options...)
	}

	planner := velty.New(options...)

	for k, v := range d.functions {
		err := planner.RegisterFunction(k, v)
		if !assert.Nil(t, err, d.description) {
			return nil, nil, err
		}
	}

	for k, v := range d.definedVars {
		err := planner.DefineVariable(k, v)
		if !assert.Nil(t, err, d.description) {
			return nil, nil, err
		}
	}

	for _, v := range d.embeddedVars {
		err := planner.EmbedVariable(v)
		if !assert.Nil(t, err, d.description) {
			return nil, nil, err
		}
	}

	exec, newState, err := planner.Compile([]byte(d.template))

	if err != nil {
		return nil, nil, err
	}

	state := newState()
	if err := d.populateState(t, state); err != nil {
		return nil, nil, err
	}

	return exec, state, nil
}

func (d *testdata) populateState(t *testing.T, state *est.State) error {
	for k, v := range d.definedVars {
		err := state.SetValue(k, v)
		if !assert.Nil(t, err, d.description+" var "+k) {
			return err
		}
	}

	for k, v := range d.embeddedVars {
		err := state.SetValue(k, v)
		if !assert.Nil(t, err, d.description+" var "+k) {
			return err
		}
	}

	return nil
}
