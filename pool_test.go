package velty

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	testcases := []struct {
		description string
		template    string
		data        map[string]interface{}
		poolSize    int
	}{
		{
			description: "poolsize 1",
			template:    "$foo ------ $var ---- $counter",
			data: map[string]interface{}{
				"foo":     "ABC",
				"var":     123,
				"counter": 0,
			},
			poolSize: 1,
		},
	}

outer:
	for _, testcase := range testcases {
		planner := New(BufferSize(1000))
		for k, v := range testcase.data {
			if !assert.Nil(t, planner.DefineVariable(k, v), testcase.description) {
				continue outer
			}
		}

		compile, newState, err := planner.Compile([]byte(testcase.template))
		if !assert.Nil(t, err, testcase.description) {
			continue
		}

		pool := NewPool(testcase.poolSize, newState)
		wg := sync.WaitGroup{}
		goroutines := 20
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(i int) {
				state := pool.State()
				defer func() {
					pool.Put(state)
					wg.Done()
				}()

				for k, v := range testcase.data {
					if !assert.Nil(t, state.SetValue(k, v), testcase.description) {
						return
					}
				}
				if !assert.Nil(t, state.SetValue("counter", i), testcase.description) {
					return
				}
				compile.Exec(state)
			}(i)
		}

		wg.Wait()

		fmt.Println(pool.counter)
	}
}
