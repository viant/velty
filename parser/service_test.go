package parser

import (
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
			description: "selector",
			input:       "${VARIABLE}",
			output:      `{"Stmt": [{"ID": "VARIABLE"}]}`,
		},
		{
			description: "empty selector",
			input:       "${}",
			expectError: true,
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
			output:      `{ "Stmt": [ { "Condition": { "X": { "X": { "X": { "Value": "1" }, "Token": "==", "Y": { "Value": "1" } }, "Token": "&&", "Y": { "X": { "Value": "2" }, "Token": "==", "Y": { "Value": "2" } } }, "Token": "&&", "Y": { "X": { "X": { "Value": "3" }, "Token": "==", "Y": { "Value": "3" } }, "Token": "||", "Y": { "X": { "Value": "4"  }, "Token": "==", "Y": { "Value": "4" } } } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement binary without token and right hand",
			input:       `#if( ${LOGGED_USER} )abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "ID": "LOGGED_USER" }, "Token": "==", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
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
			description: "if statement, different types",
			input:       `#if( 1 > "1")abc#end`,
			expectError: true,
		},
		{
			description: "if statement, boolean",
			input:       `#if( true != true)abc#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "true" }, "Token": "!=", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" } ] } } ] }`,
		},
		{
			description: "if statement, add to boolean",
			input:       `#if( false + false != true)abc#end`,
			expectError: true,
		},
		{
			description: "if statement, elseif",
			input:       `#if(true)abc#elseif("abc"=="abc")cdef#end`,
			output:      `{ "Stmt": [ { "Condition": { "X": { "Value": "true" }, "Token": "==", "Y": { "Value": "true" } }, "Body": { "Stmt": [ { "Append": "abc" }, { "Append": "cdef" } ] }, "Else": { "Condition": { "X": { "Value": "abc" }, "Token": "==", "Y": { "Value": "abc" } } } } ] }`,
		},
	}

	//for _, useCase := range useCases[len(useCases)-1:] {
	for _, useCase := range useCases {
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
