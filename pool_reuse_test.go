package velty

import (
	"testing"
)

func TestPool_ReusesStateAfterPut(t *testing.T) {
	planner := New()
	if err := planner.DefineVariable("foo", ""); err != nil {
		t.Fatalf("define variable: %v", err)
	}
	_, newState, err := planner.Compile([]byte("$foo"))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	pool := NewPool(1, newState)

	s1 := pool.State()
	pool.Put(s1)

	s2 := pool.State()
	if s1 != s2 {
		t.Fatalf("expected pool to reuse state instance")
	}
	pool.Put(s2)
}
