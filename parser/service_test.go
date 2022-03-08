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
			output:      `{"ID":"VARIABLE"}`,
		},
		{
			description: "empty selector",
			input:       "${}",
			expectError: true,
		},
		{
			description: "if statement",
			input:       `#if("1"=="1")abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"RType": "string","Value": "1"},"Y":{"RType": "string","Value": "1"}}}`,
		},
		{
			description: "if statement with left selector",
			input:       `#if(${USER_ID}=="1")abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "USER_ID"},"Y":{"RType": "string","Value": "1"}}}`,
		},
		{
			description: "if statement with both side selector",
			input:       `#if(${USER_ID}==${LOGGED_USER})abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "USER_ID"},"Y": {"ID": "LOGGED_USER"}}}`,
		},
		{
			description: "if statement with both side selector",
			input:       `#if(${USER_ID}==${LOGGED_USER})abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "USER_ID"},"Y": {"ID": "LOGGED_USER"}}}`,
		},
		{
			description: "if statement with int",
			input:       `#if(${USER_ID}==1)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "USER_ID"},"Y": {"RType": "float64", "Value": 1}}}`,
		},
		{
			description: "if statement with boolean",
			input:       `#if(${LOGGED_USER}==true)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "LOGGED_USER"},"Y": {"RType": "bool", "Value": "true"}}}`,
		},
		{
			description: "if statement with float",
			input:       `#if(${LOGGED_USER}==1.005)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "LOGGED_USER"},"Y": {"RType": "float64", "Value": "1.005"}}}`,
		},
		{
			description: "if statement with greater",
			input:       `#if(${LOGGED_USER} > 1.005)abc#end`,
			output:      `{"Condition": {"Token": ">","X": {"ID": "LOGGED_USER"},"Y": {"RType": "float64", "Value": "1.005"}}}`,
		},
		{
			description: "if statement with less",
			input:       `#if(${LOGGED_USER} < 1.005)abc#end`,
			output:      `{"Condition": {"Token": "<","X": {"ID": "LOGGED_USER"},"Y": {"RType": "float64", "Value": "1.005"}}}`,
		},
		{
			description: "if statement with less or equal",
			input:       `#if(${LOGGED_USER} <= 1.005)abc#end`,
			output:      `{"Condition": {"Token": "<=","X": {"ID": "LOGGED_USER"},"Y": {"RType": "float64", "Value": "1.005"}}}`,
		},
		{
			description: "if statement with greater or equal",
			input:       `#if(${LOGGED_USER} >= 1.005)abc#end`,
			output:      `{"Condition": {"Token": ">=","X": {"ID": "LOGGED_USER"},"Y": {"RType": "float64", "Value": "1.005"}}}`,
		},
		{
			description: "if statement with not equal",
			input:       `#if(${LOGGED_USER} != 1.005)abc#end`,
			output:      `{"Condition": {"Token": "!=","X": {"ID": "LOGGED_USER"},"Y": {"RType": "float64", "Value": "1.005"}}}`,
		},
		{
			description: "if statement with negation",
			input:       `#if(!${LOGGED_USER})abc#end`,
			output:      `{"Condition":{"Token": "!", "X": {"ID": "LOGGED_USER"}}}`,
		},
		{
			description: "if statement with AND",
			input:       `#if(!${LOGGED_USER} && ${USER_ID} == 10)abc#end`,
			output:      `{"Condition": {"Token": "&&","X": {"Token": "!", "X": {"ID": "LOGGED_USER"}},"Y": {"Token": "==", "X": {"ID": "USER_ID"}, "Y": {"Value": "10"}}}}`,
		},
		{
			description: "if statement with OR",
			input:       `#if(!${LOGGED_USER} || ${USER_ID} == 10)abc#end`,
			output:      `{"Condition": {"Token": "||","X": {"Token": "!", "X": {"ID": "LOGGED_USER"}},"Y": {"Token": "==", "X": {"ID": "USER_ID"}, "Y": {"Value": "10"}}}}`,
		},
		{
			description: "if statement with brackets ( ) #1",
			input:       `#if((1==1 && 2==2) && (3 ==3 || 4 == 4))abc#end`,
			output:      `{"Condition": {"Token": "&&","X": {"Token": "&&","X": {"Token": "==","X": {"Value": "1"},"Y": {"Value": "1"}},"Y": {"Token": "==","X": {"Value": "2"},"Y": {"Value": "2"}}},"Y": {"Token": "||","X": {"Token": "==","X": {"Value": "3"}, "Y": {"Value": "3"}},"Y": {"Token": "==","X":{"Value": "4"},"Y": {"Value": "4"}}}}}`,
		},
		{
			description: "if statement binary without token and right hand",
			input:       `#if( ${LOGGED_USER} )abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"ID": "LOGGED_USER"},"Y": {"Value": "true"}}}`,
		},
		{
			description: "if statement && token",
			input:       `#if( ${LOGGED_USER} && true )abc#end`,
			output:      `{"Condition": {"Token": "&&", "X": {"ID": "LOGGED_USER"}, "Y": {"Value": "true"}}}`,
		},
		{
			description: "if statement || token",
			input:       `#if( ${LOGGED_USER} || true )abc#end`,
			output:      `{"Condition": {"Token": "||", "X": {"ID": "LOGGED_USER"}, "Y": {"Value": "true"}}}`,
		},
		{
			description: "if statement, add equation",
			input:       `#if( 2 == 1+1)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"Value": "2"},"Y": {"Token": "+","X": {"Value": "1"},"Y": {"Value": "1"}}}}`,
		},
		{
			description: "if statement, nested add equation",
			input:       `#if( 2 == 1+1+1)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"Value": "2"},"Y": {"Token": "+","X": {"Value": "1"},"Y": {"Token": "+","X": {"Value": "1"},"Y": {"Value": "1"}}}}}`,
		},
		{
			description: "if statement, sub equation",
			input:       `#if( 0 == 1 - 1)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"Value": "0"},"Y": {"Token": "-","X": {"Value": "1"},"Y": {"Value": "1"}}}}`,
		},
		{
			description: "if statement, mul equation",
			input:       `#if( 1 == 1 * 1)abc#end`,
			output:      `{"Condition": {"Token": "==","X": {"Value": "1"},"Y": {"Token": "*","X": {"Value": "1"},"Y": {"Value": "1"}}}}`,
		},
		//{
		//	description: "if statement, quo equation",
		//	input:       `#if( 1 == 1 / 1)abc#end`,
		//	output:      `{"Condition": {"Token": "==","X": {"Value": "1"},"Y": {"Token": "/","X": {"Value": "1"},"Y": {"Value": "1"}}}}`,
		//},
		{
			description: "if statement, different types",
			input:       `#if( 1 > "1")abc#end`,
			expectError: true,
		},
		{
			description: "if statement, boolean",
			input:       `#if( true != true)abc#end`,
			output:      `{"Condition": {"Token": "!=","X": {"Value": "true"},"Y": {"Value": "true"}}}`,
		},
		{
			description: "if statement, add to boolean",
			input:       `#if( false + false != true)abc#end`,
			expectError: true,
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
