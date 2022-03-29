package utils

import (
	"github.com/viant/assertly"
	"testing"
)

func TestAppendInt(t *testing.T) {
	testcases := []struct {
		description string
		value       int
		bufferSize  int
		expected    []byte
	}{
		{
			value:      100,
			bufferSize: 5,
			expected:   []byte{'1', '0', '0'},
		},
		{
			value:      532189422,
			bufferSize: 20,
			expected:   []byte{'5', '3', '2', '1', '8', '9', '4', '2', '2'},
		},
	}

	for _, testcase := range testcases {
		buffer := make([]byte, testcase.bufferSize)
		coppied := AppendInt(buffer, int64(testcase.value), 10)
		assertly.AssertValues(t, testcase.expected, buffer[:coppied])
	}
}
