package velty_test

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty"
	"github.com/viant/velty/ast/expr"
	"github.com/viant/velty/est"
	"github.com/viant/velty/est/op"
	"github.com/viant/velty/functions"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Node struct {
	DB      *sql.DB
	Request *http.Request
}

type bar struct {
	Name string
}

func (b *bar) UpperCase() string {
	return strings.ToUpper(b.Name)
}

func (b *bar) Concat(values ...string) string {
	return strings.Join(append([]string{b.Name}, values...), " ")
}

type (
	barAggregator struct{}
	barAggregates struct {
		Names []string
	}

	sliceLength struct{}
)

func (s *sliceLength) ResultType(receiver reflect.Type, call *expr.Call) (reflect.Type, error) {
	return reflect.TypeOf(0), nil
}

func (s *sliceLength) Kind() []reflect.Kind {
	return []reflect.Kind{reflect.Slice}
}

func (s *sliceLength) Handler() interface{} {
	return func(aSlice interface{}) interface{} {
		return reflect.ValueOf(aSlice).Len()
	}
}

func (b *barAggregator) ResultType(receiver reflect.Type, call *expr.Call) (reflect.Type, error) {
	return reflect.TypeOf(&barAggregates{}), nil
}

func (b *barAggregator) Kind() []reflect.Kind {
	return []reflect.Kind{reflect.Slice}
}

func (b *barAggregator) Handler() interface{} {
	return func(data []*bar) interface{} {
		result := &barAggregates{}
		for _, datum := range data {
			result.Names = append(result.Names, datum.Name)
		}

		return result
	}
}

type ITest interface {
	TestIt(b string) string
}

type TestImpl struct {
	M string
	K int
}

func (t *TestImpl) TestIt(b string) string {
	return "my test" + b
}

func TestPlanner_Compile(t *testing.T) {
	var iTest ITest
	functions.Now = func() time.Time {
		parse, _ := time.Parse("2006-01-02 15:04:05.000000000 -0700 MST", "2014-11-12 11:45:26.000000000 +0000 UTC")
		return parse
	}

	defer func() {
		functions.Now = time.Now
	}()

	type Boo struct {
		UUID  string
		Price float64
	}

	type Foo struct {
		Name string
		ID   int
		Boo  *Boo
	}

	type Values struct {
		IntValue     int `velty:"name=Int"`
		StringValue  string
		BooleanValue bool
		Values       []string
		FloatValue   float64
		Index        int
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

	type FooNames struct {
		Names []string `velty:"names=NAMES|FOO_NAMES"`
	}

	taggedStruct := &Tagged{ID: 100}

	type Bar struct {
		Columns []struct {
			Name  string `velty:"column_name"`
			Value string `velty:"column_value"`
		}
		Id int
	}

	dep := &Department{}
	reflectStruct := reflect.StructOf([]reflect.StructField{
		{
			Name: "ID",
			Type: reflect.TypeOf(0),
		},
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
		},
		{
			Name: "Active",
			Type: reflect.TypeOf(true),
		},
	})

	type PtrsWrapper struct {
		Int     *int
		Int8    *int8
		Int16   *int16
		Int32   *int32
		Int64   *int64
		UInt    *uint
		UInt8   *uint8
		UInt16  *uint16
		UInt32  *uint32
		UInt64  *uint64
		Float64 *float64
		Float32 *float32
		String  *string
		Bool    *bool
	}

	type BoolPtrValueWrapper struct {
		ID     int
		Values *PtrsWrapper
	}

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
#if($var1 != $var2)
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
#if($var1 != $var2)
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
			description: "tags #2",
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
		{
			description: "foreach over non existing slice",
			template:    `#foreach($foo in $Foos) abc #end`,
			expect:      ``,
		},
		{
			description: `multiple fields embeded`,
			template:    `$bar.Name`,
			expect:      `bar name`,
			variables: []Variable{
				{
					Name:  "foo",
					Value: Foo{},
				},
				{
					Value: Values{},
					Embed: true,
				},
				{
					Name:  "bar",
					Value: bar{Name: "bar name"},
				},
			},
		},
		{
			description: `defined as non-pointer, indirect access`,
			template:    `$foo.Boo.Price`,
			expect:      `125.5`,
			variables: []Variable{
				{
					Name:  "bar",
					Value: bar{Name: "bar name"},
				},
				{
					Name:  "foo",
					Value: Foo{Boo: &Boo{Price: 125.5}},
				},
			},
		},
		{
			description: `unary operator`,
			template:    `#if(${x} && ${y} == "y") test #end`,
			expect:      ` test `,
			variables: []Variable{
				{
					Name:  "x",
					Value: true,
				},
				{
					Name:  "y",
					Value: "y",
				},
			},
		},
		{
			description: `slice with multiple names`,
			template:    `#foreach($foo in $FOO_NAMES) $foo #end`,
			expect:      ` abc  def  ghi `,
			variables: []Variable{
				{
					Name:  "Values",
					Value: Values{},
					Embed: true,
				},
				{
					Name:  "FooNames",
					Value: FooNames{Names: []string{"abc", "def", "ghi"}},
					Embed: true,
				},
			},
		},
		{
			description: `built in strings.ToUpper function`,
			template:    `$strings.ToUpper("abc")`,
			expect:      `ABC`,
		},
		{
			description: `built in slices.Length function`,
			template:    `$slices.Length($foos)`,
			expect:      `3`,
			definedVars: map[string]interface{}{
				"foos": []int{1, 2, 3},
			},
		},
		{
			description: `foreach over interfaces`,
			template:    `#foreach($value in $values) $value ,#end`,
			definedVars: map[string]interface{}{
				"values": []interface{}{1, "abc", true, 2.0},
			},
			expect: ` 1 , abc , true , 2 ,`,
		},
		{
			description: `foreach over interfaces`,
			template:    `$values`,
			expect:      `[1,"abc",true,2]`,
			definedVars: map[string]interface{}{
				"values": []interface{}{1, "abc", true, 2.0},
			},
		},
		{
			description: `foreach over interfaces`,
			template:    `#foreach($value in $values) $types.IsString($value)#end`,
			definedVars: map[string]interface{}{
				"values": []interface{}{1, "abc", true, 2.0},
			},
			expect: ` false true false false`,
		},
		{
			description: `register error`,
			template:    `$errors.RegisterError("error was thrown " + $message)`,
			definedVars: map[string]interface{}{
				"message": "this is error message",
			},
			expectTemplateErr: true,
		},
		{
			description: `String in`,
			template:    `$strings.In($Set, "2")`,
			definedVars: map[string]interface{}{
				"Set": []string{"1", "2", "3"},
			},
			expect: "true",
		},
		{
			description: `assign struct`,
			template:    `#foreach($barColumn in $bar.Columns)#set($lastColumn = $barColumn)#end $lastColumn.column_name  $lastColumn.column_value`,
			definedVars: map[string]interface{}{
				"bar": Bar{Columns: []struct {
					Name  string `velty:"column_name"`
					Value string `velty:"column_value"`
				}{
					{
						Name:  "Col - 1",
						Value: "Val - 1",
					},
					{
						Name:  "Col - 2",
						Value: "Val - 2",
					},
				}},
			},
			expect: ` Col - 2  Val - 2`,
		},
		{
			description: `foreach`,
			template:    `#foreach($anInt in $ints)$anInt,#end#foreach($anInt in $ints)$anInt,#end`,
			expect:      `1,2,3,4,1,2,3,4,`,
			variables: []Variable{
				{
					Value: []int{1, 2, 3, 4},
					Name:  "ints",
				},
			},
		},
		{
			description: `or`,
			template:    `#if(("abc" == "abc") && (1 == 1))equal#elsenot equal#end`,
			expect:      `equal`,
		},
		{
			description: `set function call`,
			template: `#foreach($foo in $Foos)
#set($aVar = $strings.ToLower($foo.Name))$aVar
#end`,
			variables: []Variable{{
				Name: "Foos",
				Value: []Foo{
					{
						Name: "ABCDEF",
					},
					{
						Name: "GHI",
					},
				},
			}},
			expect: "\nabcdef\n\nghi\n",
		},
		{
			description: `chaining methods`,
			template: `#foreach($barColumn in $bar.Columns)
#set($lastColumnName = $slices.StringAt($strings.Split($barColumn.column_name, "-"), 0))
#end
$lastColumnName`,
			definedVars: map[string]interface{}{
				"bar": Bar{Columns: []struct {
					Name  string `velty:"column_name"`
					Value string `velty:"column_value"`
				}{
					{
						Name:  "Col - 1",
						Value: "Val - 1",
					},
					{
						Name:  "Col - 2",
						Value: "Val - 2",
					},
				}},
			},
			expect: "\n\n\n\n\nCol ",
		},
		{
			description: `reverse strings`,
			template: ` 
		#foreach($column in $foo.columns)
			#foreach($value in $column.values)
					#if($errors.AssertString($value))
						#set($reversed = $slices.ReverseStrings($strings.Split($value, ",")))
						#foreach($value in $reversed)$value#end
					#end
			#end
		#end
	`,
			definedVars: map[string]interface{}{
				"foo": filter{
					columns: []column{
						{
							values: []interface{}{"abc,def", "ghi,jkl"},
						},
					},
				},
			},
			expect: " \n\t\t\n\t\t\t\n\t\t\t\t\t\n\t\t\t\t\t\t\n\t\t\t\t\t\tdefabc\n\t\t\t\t\t\n\t\t\t\n\t\t\t\t\t\n\t\t\t\t\t\t\n\t\t\t\t\t\tjklghi\n\t\t\t\t\t\n\t\t\t\n\t\t\n\t",
		},
		{
			description: `time now`,
			template:    `$time.Now()`,
			expect:      `"2014-11-12T11:45:26Z"`,
		},
		{
			description: `slices`,
			template:    `$values[2]`,
			expect:      `second value`,
			definedVars: map[string]interface{}{
				"values": []string{"zero value", "first value", "second value", "third value"},
			},
		},
		{
			description: `types difference`,
			template:    `#set($value = 2.5 + $anInt)$value`,
			expect:      `3.5`,
			definedVars: map[string]interface{}{
				"anInt": 1,
			},
		},
		{
			description: `types difference #2`,
			template:    `#set($value = ("Values: " + 1) + (" another one: " + 5.21))$value`,
			expect:      `Values: 1 another one: 5.21`,
		},
		{
			description: `map get | exported`,
			template:    `$Data["key"]`,
			expect:      `value`,
			definedVars: map[string]interface{}{
				"Data": map[string]interface{}{
					"key": "value",
				},
			},
		},
		{
			description: `map get | unexported`,
			template:    `$data["key"]`,
			expect:      `value`,
			definedVars: map[string]interface{}{
				"data": map[string]interface{}{
					"key": "value",
				},
			},
		},
		{
			description: `map get | has true`,
			template:    `$data.HasKey("key")`,
			expect:      `true`,
			definedVars: map[string]interface{}{
				"data": map[string]interface{}{
					"key": "value",
				},
			},
		},
		{
			description: `map get | has false`,
			template:    `$data.HasKey("abc")`,
			expect:      `false`,
			definedVars: map[string]interface{}{
				"data": map[string]interface{}{
					"key": "value",
				},
			},
		},
		{
			description: `map get | has false`,
			template:    `$data.HasKey("abc")`,
			expect:      `false`,
			definedVars: map[string]interface{}{
				"data": map[string]interface{}{
					"key": "value",
				},
			},
		},
		{
			description: `slice | index by`,
			template:    `#set($index = $Values.IndexBy("StringValue")) $index.HasKey("key - 1")  $index.HasKey("")  $index.HasKey("key - 3")`,
			expect:      ` true  false  true`,
			definedVars: map[string]interface{}{
				"Values": []*Values{
					{
						StringValue: "key - 1",
					},
					{
						StringValue: "key - 2",
					},
					{
						StringValue: "key - 3",
					},
				},
			},
		},
		{
			description:       `slice | index by`,
			template:          `$strconv.Atoi("adewqdsa") abcdef`,
			expect:            ``,
			options:           []velty.Option{velty.PanicOnError(true)},
			expectTemplateErr: true,
		},
		{
			description: `slice | index by, type missmatch`,
			template:    `#set($index = $Values.IndexBy("IntValue")) $index[$aKey].StringValue`,
			expect:      ` key - 1`,
			definedVars: map[string]interface{}{
				"aKey": int16(123),
				"Values": []*Values{
					{
						StringValue: "key - 1",
						IntValue:    123,
					},
					{
						StringValue: "key - 2",
						IntValue:    246,
					},
					{
						StringValue: "key - 3",
						IntValue:    492,
					},
				},
			},
		},
		{
			description: `multimap | single value`,
			template:    `$values[1][1]`,
			expect:      `10`,
			definedVars: map[string]interface{}{
				"values": map[int]map[int]int{
					1: {
						1: 10,
					},
				},
			},
		},
		{
			description: `multimap | multivalues`,
			template:    `$values[2][3]`,
			expect:      `30`,
			definedVars: map[string]interface{}{
				"values": map[int]map[int]int{
					1: {
						1: 10,
						2: 20,
					},
					2: {
						3: 30,
						4: 40,
					},
				},
			},
		},
		{
			description: `multimap | multivalues`,
			template:    `${values[2][3]}`,
			expect:      `30`,
			definedVars: map[string]interface{}{
				"values": map[interface{}]interface{}{
					1: map[interface{}]interface{}{
						1: 10,
						2: 20,
					},
					2: map[interface{}]interface{}{
						3: 30,
						4: 40,
					},
				},
			},
		},
		{
			description: "autogen anonymous T",
			template:    "$ID | $Name | $Active",
			expect:      `0 |  | false`,
			embeddedVars: map[string]interface{}{
				"": reflect.New(reflectStruct).Elem().Interface(),
			},
		},
		{
			description: "autogen anonymous *T",
			template:    "$ID | $Name | $Active",
			expect:      `0 |  | false`,
			embeddedVars: map[string]interface{}{
				"": reflect.New(reflectStruct).Interface(),
			},
		},
		{
			description: "modify given struct",
			template:    `#set($Foo.Name = "changed")$Foo.Name`,
			expect:      `changed`,
			definedVars: map[string]interface{}{
				"Foo": &Foo{},
			},
		},
		{
			description: `db | cyclic deps, interfaces`,
			template:    `123`,
			expect:      `123`,
			definedVars: map[string]interface{}{
				"DB": &sql.DB{},
			},
		},
		{
			description: `node | cyclic deps`,
			template:    `123`,
			expect:      `123`,
			definedVars: map[string]interface{}{
				"DB": &Node{},
			},
		},
		{
			description: "assign pointers",
			template:    "#set($assigned = $ptrValue)$assigned",
			expect:      `10`,
			definedVars: map[string]interface{}{
				"ptrValue": intPtr(10),
			},
		},
		{
			description: "indirect assign",
			template: `
			#set($ptrWrapper.Int = ${int8})
			$ptrWrapper.Int
`,
			expect: `
			
			1
`,
			definedVars: map[string]interface{}{
				"ptrWrapper": &PtrsWrapper{},
				"int8":       int8(1),
			},
		},
		{
			description: "result typer",
			template:    `#set($Aggregated = $bars.Aggregates()) ${Aggregated.Names.Length()}`,
			expect:      ` 3`,
			definedVars: map[string]interface{}{
				"bars": []*bar{{Name: "abc"}, {Name: "def"}, {Name: "ghi"}},
			},
			KindFunctions: map[string]op.KindFunction{
				"Aggregates": &barAggregator{},
				"Length":     &sliceLength{},
			},
		},
		{
			description: "register type func",
			template:    `$sizeable.Size()`,
			expect:      `4`,
			definedVars: map[string]interface{}{
				"sizeable": []int{1, 2, 3, 4},
			},
			typeFunc: map[reflect.Type][]*op.TypeFunc{
				reflect.TypeOf([]int{}): {
					{
						Name: "Size",
						Handler: func(data []int) int {
							return len(data)
						},
					},
				},
			},
		},
		{
			description: "iterate over nil slice",
			template:    `#foreach($item in $Foos.Items) $item #end`,
			expect:      ``,
			definedVars: map[string]interface{}{
				"Foos": &struct {
					Items []string
					_tmp  int
				}{},
			},
		},
		{
			description: `embeded pointer`,
			template:    `${StringValue} |#foreach($aStr in $Values) ${aStr} #end`,
			expect:      `a value | 1  2  3 `,
			embeddedVars: map[string]interface{}{
				"bar": &Values{
					StringValue: "a value",
					Values:      []string{"1", "2", "3"},
				},
			},
		},
		{
			description: `if struct not nil`,
			template:    `#if($foo)foo was set#else unset foo#end`,
			expect:      "foo was set",
			definedVars: map[string]interface{}{
				"foo": &Foo{},
			},
		},
		{
			description: `if struct not nil`,
			template:    `#if($foo)foo was set#else unset foo#end`,
			expect:      " unset foo",
			definedVars: map[string]interface{}{
				"foo": dep.Address,
			},
		},
		{
			description: `json`,
			template:    `#set($aFoo = $json.UnmarshalInto($fooBody, "*Foo") )$aFoo.ID | $aFoo.Name | $aFoo.Price`,
			expect:      `10 | Foo | 125.5`,
			definedVars: map[string]interface{}{
				"fooBody": `{"ID": 10, "Name": "Foo", "Price": 125.5}`,
			},
			options: []velty.Option{
				velty.TypeParser(func(typeRepresentation string) (reflect.Type, error) {
					type Foo struct {
						ID    int
						Name  string
						Price float64
					}

					return reflect.TypeOf(&Foo{}), nil
				}),
			},
		},
		{
			template: `$ToUpperCase("value")`,
			standaloneFunctions: map[string]*op.Function{
				"ToUpperCase": {
					ResultTyper: func(call *expr.Call) (reflect.Type, error) {
						return reflect.TypeOf(""), nil
					},
					Handler: func(value string) interface{} {
						return strings.ToUpper(value)
					},
				},
			},
			expect: "VALUE",
		},
		{
			template: `#set($Values = $wrapper.Values)#if($Values.Bool)is true#elseis false#end`,
			definedVars: map[string]interface{}{
				"wrapper": &BoolPtrValueWrapper{
					Values: &PtrsWrapper{
						Bool: boolPtr(true),
					},
				},
			},
			expect: "is true",
		},

		{
			description: "iface call",
			template:    `$tester.TestIt("true")`,
			definedVars: map[string]interface{}{
				"tester": &definedVariable{
					valueType: reflect.TypeOf(&iTest).Elem(),
					value:     &TestImpl{M: "abc"},
				},
			},
			expect: "my testtrue",
		},
	}

	//for i, testCase := range testCases[:len(testCases)-1] {
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

		err = exec.Exec(state)
		if testCase.expectTemplateErr {
			assert.NotNil(t, err, testCase.description)
		} else {
			assert.Nil(t, err, testCase.description)
		}
		output := state.Buffer.Bytes()
		assert.Equal(t, testCase.expect, string(output), testCase.description)
	}
}

type testdata struct {
	description         string
	template            string
	definedVars         map[string]interface{}
	embeddedVars        map[string]interface{}
	functions           map[string]interface{}
	variables           []Variable
	expectError         bool
	expect              string
	options             []velty.Option
	expectTemplateErr   bool
	setVariables        map[string]interface{}
	KindFunctions       map[string]op.KindFunction
	typeFunc            map[reflect.Type][]*op.TypeFunc
	standaloneFunctions map[string]*op.Function
}

type Variable struct {
	Name  string
	Value interface{}
	Embed bool
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
		if def, ok := v.(*definedVariable); ok {
			planner.DefineVariable(k, def.valueType)
			continue
		}

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

	for i, variable := range d.variables {
		if variable.Embed {
			err := planner.EmbedVariable(d.variables[i].Value)
			if !assert.Nil(t, err, d.description) {
				return nil, nil, err
			}
		} else {
			err := planner.DefineVariable(variable.Name, d.variables[i].Value)
			if !assert.Nil(t, err, d.description) {
				return nil, nil, err
			}
		}
	}

	for name, function := range d.KindFunctions {
		if err := planner.RegisterFunctionKind(name, function); err != nil {
			return nil, nil, err
		}
	}

	for rType, funcs := range d.typeFunc {
		for _, aFunc := range funcs {
			if err := planner.RegisterTypeFunc(rType, aFunc); err != nil {
				return nil, nil, err
			}
		}
	}

	for fnName, functioner := range d.standaloneFunctions {
		if err := planner.RegisterStandaloneFunction(fnName, functioner); err != nil {
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
		if def, ok := v.(*definedVariable); ok {
			v = def.value
		}

		err := state.SetValue(k, v)
		if !assert.Nil(t, err, d.description+" var "+k) {
			return err
		}
	}

	for k, v := range d.embeddedVars {
		err := state.EmbedValue(v)
		if !assert.Nil(t, err, d.description+" var "+k) {
			return err
		}
	}

	for i, variable := range d.variables {
		if variable.Embed {
			err := state.EmbedValue(d.variables[i].Value)
			if !assert.Nil(t, err, d.description) {
				return err
			}
		} else {
			err := state.SetValue(variable.Name, d.variables[i].Value)
			if !assert.Nil(t, err, d.description) {
				return err
			}
		}
	}

	for key := range d.setVariables {
		if err := state.SetValue(key, d.setVariables[key]); err != nil {
			return err
		}
	}

	return nil
}

func Test_ForEach_Issue(t *testing.T) {

	type Foo struct {
		Bar  string `velty:"names=Bar"`
		Bar1 int    `velty:"names=Bar1"`
	}

	type Repeated struct {
		URLs []string `velty:"names=URLS"`
	}
	type Test struct {
		Foo // when foo is commented out it's working otherwise it does not, address for slice not correctly computed
		Repeated
	}

	tmpl := `#foreach ($url in ${URLS})
<img src="${url}" style="display:none" height="1" width="1">
#end`
	planner := velty.New()
	planner.EmbedVariable(Test{})
	exec, newState, err := planner.Compile([]byte(tmpl))
	assert.Nil(t, err)
	aState := newState()
	aTest := Test{}
	aTest.URLs = []string{"urtl1", "urtl2", "urtl3"}
	aState.EmbedValue(aTest)

	assert.Nil(t, exec.Exec(aState))
	actual := aState.Buffer.String()

	expect := `
<img src="urtl1" style="display:none" height="1" width="1">

<img src="urtl2" style="display:none" height="1" width="1">

<img src="urtl3" style="display:none" height="1" width="1">
`
	assert.EqualValues(t, expect, actual)

}

type definedVariable struct {
	valueType interface{}
	value     interface{}
}

type filter struct {
	columns []column
}

type column struct {
	values  []interface{}
	aBool   bool
	aString string
}

func intPtr(value int) *int {
	return &value
}

func boolPtr(aBool bool) *bool {
	return &aBool
}
