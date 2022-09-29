package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"testing"
)

func TestService_Parse(t *testing.T) {
	useCases := []struct {
		input       string
		output      string
		description string
		expectError bool
	}{
		{
			description: "without velocity tags",
			input:       "<h3>Some h3 title</h3>",
			output:      `{ "Stmt": [ { "Append": "<h3>Some h3 title</h3>" } ] }`,
		},
		{
			description: "selector",
			input:       "${VARIABLE}",
			output:      `{"Stmt": [{"ID": "VARIABLE"}]}`,
		},
		{
			description: "if statement",
			input:       `#if("1"=="1")abc#end`,
			output:      `{"Stmt": [{"Condition": {"X": {"Value": "1"}, "Token": "==", "Y": {"Value": "1"}}, "Body": {"Stmt": [{"Append": "abc"}]}}]}`,
		},
		{
			description: "if statement with left selector",
			input:       `#if(${USER_ID}=="1")abc#end`,
			output:      `{"Stmt": [{"Condition": {"X": { "ID": "USER_ID" }, "Token": "==", "Y": { "Value": "1" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with both side selector",
			input:       `#if(${USER_ID}==${LOGGED_USER})abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "USER_ID" }, "Token": "==", "Y": { "ID": "LOGGED_USER" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with int",
			input:       `#if(${USER_ID}==1)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "USER_ID" }, "Token": "==", "Y": { "Value": "1" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with boolean",
			input:       `#if(${LOGGED_USER}==true)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "==", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with number",
			input:       `#if(${LOGGED_USER}==1.005)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "==", "Y": { "Value": "1.005" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with greater",
			input:       `#if(${LOGGED_USER} > 1.005)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": ">", "Y": { "Value": "1.005" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with less",
			input:       `#if(${LOGGED_USER} < 1.005)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "<", "Y": { "Value": "1.005" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with less or equal",
			input:       `#if(${LOGGED_USER} <= 1.005)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "<=", "Y": { "Value": "1.005" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with greater or equal",
			input:       `#if(${LOGGED_USER} >= 1.005)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": ">=", "Y": { "Value": "1.005" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with not equal",
			input:       `#if(${LOGGED_USER} != 1.005)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "!=", "Y": { "Value": "1.005" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with negation",
			input:       `#if(!${LOGGED_USER})abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "Token": "!", "X": { "ID": "LOGGED_USER" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with AND",
			input:       `#if(!${LOGGED_USER} && ${USER_ID} == 10)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Token": "!", "X": { "ID": "LOGGED_USER" } }, "Token": "&&", "Y": { "X": { "ID": "USER_ID" }, "Token": "==", "Y": { "Value": "10" } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with OR",
			input:       `#if(!${LOGGED_USER} || ${USER_ID} == 10)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Token": "!", "X": { "ID": "LOGGED_USER" } }, "Token": "||", "Y": { "X": { "ID": "USER_ID" }, "Token": "==", "Y": { "Value": "10" } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement with brackets ( ) #1",
			input:       `#if((1==1 && 2==2) && (3 ==3 || 4 == 4))abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "P": { "X": { "X": { "Value": "1" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "&&", "Y": { "X": { "Value": "2" }, "Token": "==", "Y": { "Value": "2" } } } }, "Token": "&&", "Y": { "P": { "X": { "Value": "3" }, "Token": "==", "Y": { "X": { "Value": "3" }, "Token": "||", "Y": { "X": { "Value": "4" }, "Token": "==", "Y": { "Value": "4" } } } } } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement binary without token and right hand",
			input:       `#if( ${LOGGED_USER} )abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "ID": "LOGGED_USER" }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement && token",
			input:       `#if( ${LOGGED_USER} && true )abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "&&", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement || token",
			input:       `#if( ${LOGGED_USER} || true )abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "||", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, add equation",
			input:       `#if( 2 == 1+1)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "2" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "+", "Y": { "Value": "1" } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, nested add equation",
			input:       `#if( 2 == 1+1+1)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "2" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "+", "Y": { "X": { "Value": "1" }, "Token": "+", "Y": { "Value": "1" } } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, sub equation",
			input:       `#if( 0 == 1 - 1)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "0" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "-", "Y": { "Value": "1" } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, mul equation",
			input:       `#if( 1 == 1 * 1)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "1" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "*", "Y": { "Value": "1" } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, quo equation",
			input:       `#if( 1 == 1 / 1)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "1" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "/", "Y": { "Value": "1" } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, boolean",
			input:       `#if( true != true)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "true" }, "Token": "!=", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, elseif",
			input:       `#if(true)abc#elseif("abc"=="abc")cdef#end`,
			output:      `{ "Stmt": [ { "Condition": { "Value": "true" }, "Body": { "Stmt": [ { "Append": "abc" } ] }, "Else": { "Condition": { "X": { "Value": "abc" }, "Token": "==", "Y": { "Value": "abc" } }, "Body": { "Stmt": [ { "Append": "cdef" } ] } } } ] }`,
		},
		{
			description: "if statement, else",
			input:       `#if(true)abc#elsecdef#end`,
			output:      `{ "Stmt": [ { "Condition": { "Value": "true" }, "Body": { "Stmt": [ { "Append": "abc" } ] }, "Else": { "Condition": { "X": { "Value": "true" }, "Token": "==", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "cdef" } ] } } } ] }`,
		},
		{
			description: "set value",
			input:       `#set ($message="Hello World")`,
			output:      `{ "Stmt": [ { "X": { "ID": "message" }, "Op": "=", "Y": { "Value": "Hello World" } } ] }`,
		},
		{
			description: "set value as equation",
			input:       `#set ($value= 10 + 10 * 10)`,
			output:      `{ "Stmt": [ { "X": { "ID": "value" }, "Op": "=", "Y": { "X": { "Value": "10" }, "Token": "+", "Y": { "X": { "Value": "10" }, "Token": "*", "Y": { "Value": "10" } } } } ] }`,
		},
		{
			description: "foreach",
			input:       `<ul>#foreach( $value in $values)<li>${value}</li>#end</ul>`,
			output:      `{ "Stmt": [ { "Append": "<ul>" }, { "Item": { "ID": "value" }, "Set": { "ID": "values" }, "Body": { "Stmt": [ { "Append": "<li>" }, { "ID": "value" }, { "Append": "</li>" } ] } }, { "Append": "</ul>" } ] }`,
		},
		{
			description: "foreach with index",
			input:       `<ul>#foreach( $value, $index in $values)<li>${value}, ${index}</li>#end</ul>`,
			output:      `{ "Stmt": [ { "Append": "<ul>" }, { "Index": { "ID": "index" }, "Item": { "ID": "value" }, "Set": { "ID": "values" }, "Body": { "Stmt": [ { "Append": "<li>" }, { "ID": "value" }, { "Append": ", " }, { "ID": "index" }, { "Append": "</li>" } ] } }, { "Append": "</ul>" } ] }`,
		},
		{
			description: "for loop",
			input:       `<ul>#for( $var = 1; $var < 10; $var++ )<li>${value}, ${index}</li>#end</ul>`,
			output:      `{ "Stmt": [ { "Append": "<ul>" }, { "Init": { "X": { "ID": "var" }, "Op": "=", "Y": { "Value": "1" } }, "Cond": { "X": { "ID": "var" }, "Token": "<", "Y": { "Value": "10" } }, "Body": { "Stmt": [ { "Append": "<li>" }, { "ID": "value" }, { "Append": ", " }, { "ID": "index" }, { "Append": "</li>" } ] }, "Post": { "X": { "ID": "var" }, "Op": "=", "Y": { "X": { "ID": "var" }, "Token": "+", "Y": { "Value": "1" } } } }, { "Append": "</ul>" } ] }`,
		},
		{
			description: "different selectors",
			input:       `#if( $id == ${Id3.Name} )#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "id" }, "Token": "==", "Y": { "ID": "Id3", "X": { "ID": "Name" } } } } ] }`,
		},
		{
			description: "selector without brackets and number",
			input:       `#if( $id == 1 )#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "id" }, "Token": "==", "Y": { "Value": "1" } } } ] }`,
		},
		{
			description: "multiple comparisons",
			input:       `#if( $id == 1 == true == false )#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "id" }, "Token": "==", "Y": { "X": { "Value": "1" }, "Token": "==", "Y": { "X": { "Value": "true" }, "Token": "==", "Y": { "Value": "false" } } } } } ] }`,
		},
		{
			description: "function call",
			input:       `${foo.Function(123, !true, -5, 123+321, 10 * 10,!${USER_LOGGED})}`,
			output:      `{ "Stmt": [ { "ID": "foo", "X": { "ID": "Function", "X": { "Args": [ { "Value": "123" }, { "Token": "!", "X": { "Value": "true" } }, { "Value": "-5" }, { "X": { "Value": "123" }, "Token": "+", "Y": { "Value": "321" } }, { "X": { "Value": "10" }, "Token": "*", "Y": { "Value": "10" } }, { "Token": "!", "X": { "ID": "USER_LOGGED" } } ] } } } ] }`,
		},
		{
			description: "comments",
			input: `## THIS IS COMMENT
#if(1==1)abc#end
## THIS IS ALSO COMMENT`,
			output: `{ "Stmt": [ { "Condition": { "X": { "Value": "1" }, "Token": "==", "Y": { "Value": "1" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: `assign binary expression`,
			input:       `#set( $var1 = $foo + 10)abc`,
			output:      `{ "Stmt": [ { "X": { "ID": "var1" }, "Op": "=", "Y": { "X": { "ID": "foo" }, "Token": "+", "Y": { "Value": "10" } } }, { "Append": "abc" } ] }`,
		},
		{
			description: `assign binary expression`,
			input:       `#set( $var1 = $foo != 10)abc`,
			output:      `{ "Stmt": [ { "X": { "ID": "var1" }, "Op": "=", "Y": { "X": { "ID": "foo" }, "Token": "!=", "Y": { "Value": "10" } } }, { "Append": "abc" } ] }`,
		},
		{
			description: `evaluate`,
			input:       `#evaluate(${FOO_TEMPLATE})`,
			output:      `{ "Stmt": [ { "X": { "ID": "FOO_TEMPLATE" } } ] }`,
		},
		{
			description: `selector without brackets`,
			input:       `$FOO.VALUES.NAME<h3>`,
			output:      `{ "Stmt": [ { "ID": "FOO", "X": { "ID": "VALUES", "X": { "ID": "NAME" } } }, { "Append": "<h3>" } ] }`,
		},
		{
			description: `evaluate`,
			input:       `$!FOO.VALUES.NAME<h3>`,
			output:      `{ "Stmt": [ { "ID": "FOO", "X": { "ID": "VALUES", "X": { "ID": "NAME" } } }, { "Append": "<h3>" } ] }`,
		},
		{
			description: `range`,
			input:       `#foreach($int in [-10...10]) abc #end`,
			output:      `{ "Stmt": [ { "Item": { "ID": "int" }, "Set": { "X": { "Value": "-10" }, "Y": { "Value": "10" } }, "Body": { "Stmt": [ { "Append": " abc " } ] } } ] }`,
		},
		{
			description: `method call`,
			input:       `$bar.Concat($foo, $var.toUpperCase(), "abcdef")`,
			output:      `{ "Stmt": [ { "ID": "bar", "X": { "ID": "Concat", "X": { "Args": [ { "ID": "foo" }, { "ID": "var", "X": { "ID": "toUpperCase", "X": { "Args": [] } } }, { "Value": "abcdef" } ] } } } ] }`,
		},
		{
			description: `empty input`,
			input:       ``,
			output:      `{"Stmt":[]}`,
		},
		{
			description: `stmt block`,
			input:       `#if(1==1)abc#{else}def#{end}`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "1" }, "Token": "==", "Y": { "Value": "1" } }, "Body": { "Stmt": [ { "Append": "abc" } ] }, "Else": { "Condition": { "X": { "Value": "true" }, "Token": "==", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "def" } ] } } } ] }`,
		},
		{
			description: `stmt block`,
			input:       `http://localhost:8080/index#123`,
			output:      `{ "Stmt": [ { "Append": "http://localhost:8080/index" }, { "Append": "#" }, { "Append": "123" } ] }`,
		},
		{
			description: "!$",
			input:       `#if(${abc} && !("$!{def}"=="")) abc #end`,
			output:      `{ "Stmt": [ { "Body": { "Stmt": [ { "Append": " abc " } ] }, "Condition": { "Token": "&&", "X": { "FullName": "${abc}", "ID": "abc" }, "Y": { "Token": "!", "X": { "P": { "Token": "==", "X": { "FullName": "${def}", "ID": "def" }, "Y": { "Value": "" } } } } } } ] }`,
		},
		{
			description: "$ before selector",
			input:       `$${abc}`,
			output:      `{ "Stmt": [ { "Append": "$" }, { "ID": "abc", "FullName": "${abc}" } ] }`,
		},
		{
			description: `for with !$`,
			input:       ` #foreach ($abc in $!{Collection}) forEach body #end`,
			output:      `{ "Stmt": [ { "Append": " " }, { "Item": { "ID": "abc", "FullName": "" }, "Set": { "ID": "Collection", "FullName": "${Collection}" }, "Body": { "Stmt": [ { "Append": " forEach body " } ] } } ] }`,
		},
		{
			description: `stmt block #2`,
			input: `#if (${foo} != 1)
            #if(${boo}==2)abc#{else}def#{end}?';
#else
#end
`,
			output: `{"Stmt": [ { "Condition": { "X": { "ID": "foo", "FullName": "${foo}" }, "Token": "!=", "Y": { "Value": "1" } }, "Body": { "Stmt": [ { "Append": "\n            " }, { "Condition": { "X": { "ID": "boo", "FullName": "${boo}" }, "Token": "==", "Y": { "Value": "2" } }, "Body": { "Stmt": [ { "Append": "abc" } ] }, "Else": { "Condition": { "X": { "Value": "true" }, "Token": "==", "Y": { "Value": "true" }  },  "Body": { "Stmt": [ { "Append": "def"  } ] } } }, { "Append": "?';\n" } ] }, "Else": { "Condition": { "X": { "Value": "true" }, "Token": "==", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "\n" } ] } } }, { "Append": "\n" } ] }`,
		},
		{
			description: `escape selector`,
			input:       `#[[$ABC_DEF]]#.Columns[$i].Values[$j]`,
			output:      `{ "Stmt": [ { "Append": "$ABC_DEF" }, { "Append": ".Columns[" }, { "ID": "i", "FullName": "$i" }, { "Append": "].Values[" }, { "ID": "j", "FullName": "$j" }, { "Append": "]" } ] }`,
		},
		{
			description: `slices chain`,
			input:       `$Columns[$i].Values[$j].Value`,
			output:      ` { "Stmt": [ { "ID": "Columns", "X": { "X": { "ID": "i", "FullName": "$i" }, "Y": { "ID": "Values", "X": { "X": { "ID": "j", "FullName": "$j" }, "Y": { "ID": "Value", "FullName": "$Value" } }, "FullName": "$Values[$j].Value" } }, "FullName": "$Columns[$i].Values[$j].Value" } ] }`,
		},
		{
			description: `assign in If`,
			input:       `#set($aValue = 1 ) #if($aValue = 1) abc #end def`,
			expectError: true,
		},
		{
			description: `multiple parentheses`,
			input:       `#set($value = ("Values: " + 1) + (" another one: " + 5.21))$value`,
			output:      `{ "Stmt": [ { "X": { "ID": "value", "FullName": "" }, "Op": "=", "Y": { "X": { "P": { "X": { "Value": "Values: " }, "Token": "+", "Y": { "Value": "1" } } }, "Token": "+", "Y": { "P": { "X": { "Value": " another one: " }, "Token": "+", "Y": { "Value": "5.21" } } } } }, { "ID": "value", "FullName": "$value" } ] }`,
		},
	}

	//for i, useCase := range useCases[len(useCases)-1:] {
	for i, useCase := range useCases {
		fmt.Printf("Running testcase: %v\n", i)
		node, err := Parse([]byte(useCase.input))

		if useCase.expectError {
			assert.NotNil(t, err, useCase.description)
			continue
		}

		assert.Nil(t, err, useCase.description)
		if !assertly.AssertValues(t, useCase.output, node, useCase.description) {
			toolbox.DumpIndent(node, true)
		}
	}
}
