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
			template: `
			#set($percent = $number / 100)$percent`,
			vars: map[string]interface{}{
				"number": 12000,
			},
			expect: `120`,
		},
	}

	for _, testCase := range testCases {
		planner := plan.New()

		for k, v := range testCase.vars {
			err := planner.DefineVariable(k, v)
			if !assert.Nil(t, err, testCase.description) {
				continue
			}
		}
		exec, newState, err := planner.Compile(testCase.template)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		state := newState()
		for k, v := range testCase.vars {
			state.SetValue(k, v)
		}
		exec.Exec(state)
		output := state.Buffer.Bytes()
		assert.Equal(t, testCase.expect, string(output), testCase.description)
	}

}
