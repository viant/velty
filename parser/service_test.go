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
