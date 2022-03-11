package plan_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/velty/est/plan"
	"testing"
)

func TestPlanner_Compile(t *testing.T) {

	var testCases = []struct {
		description string
		template    string
		vars        map[string]interface{}
		expect      string
	}{
		{
			template: `#set($var1 = 123)$var1`,
			expect:   `123`,
		},
		{
			template: `#set($var1 = 12400/100)$var1`,
			expect:   `124`,
		},
		{
			template: `#set($var1 = 3 + $num)$var1`,
			expect:   `13`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			template: `#set($var1 = $num +  3)$var1`,
			expect:   `13`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
		{
			template: `#set($var1 = $num - 3)$var1`,
			expect:   `7`,
			vars: map[string]interface{}{
				"num": 10,
			},
		},
	}
outer:
	for _, testCase := range testCases {
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
