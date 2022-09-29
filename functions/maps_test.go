package functions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaps_Has(t *testing.T) {
	testCases := []struct {
		description string
		aMap        interface{}
		key         interface{}
		expectFound bool
		expectErr   bool
		value       interface{}
	}{
		{
			aMap: map[string]int{
				"abc": 1,
			},
			key:         "abc",
			expectFound: true,
			value:       int(1),
		},
	}

	maps := Maps{}
	for _, testCase := range testCases {
		has, err := maps.Has(testCase.aMap, testCase.key)
		if !testCase.expectErr && !assert.Nil(t, err, testCase.description) {
			continue
		}

		assert.Equal(t, testCase.expectFound, has)
		get, err := maps.Get(testCase.aMap, testCase.key)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		assert.Equal(t, testCase.value, get, testCase.description)
	}
}
