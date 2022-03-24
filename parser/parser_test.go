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
			input:       `#if( $id == ${id3.Name} )#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "id" }, "Token": "==", "Y": { "ID": "id3", "X": { "ID": "Name" } } } } ] }`,
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
			output:      `{ "Stmt": [ { "ID": "foo", "X": { "ID": "Function", "X": { "Args": [ { "Value": "123" }, { "Value": "123" }, { "Token": "!", "X": { "Value": "true" } }, { "Value": "123" }, { "Token": "-", "X": { "Value": "5" } }, { "Value": "123" }, { "X": { "Value": "123" }, "Token": "+", "Y": { "Value": "321" } }, { "Value": "123" }, { "X": { "Value": "10" }, "Token": "*", "Y": { "Value": "10" } }, { "Value": "123" }, { "Token": "!", "X": { "ID": "USER_LOGGED" } } ] } } } ] }`,
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
