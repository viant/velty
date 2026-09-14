package velty

import (
	"context"
	"errors"
	"testing"
)

func TestPool_ReturnsUsableStateAfterPut(t *testing.T) {
	planner := New()
	if err := planner.DefineVariable("foo", ""); err != nil {
		t.Fatal(err)
	}
	executable, newState, err := planner.Compile([]byte("$foo"))
	if err != nil {
		t.Fatal(err)
	}
	pool := NewPool(1, newState)
	first := pool.State()
	if err = first.SetValue("foo", "first"); err != nil {
		t.Fatal(err)
	}
	if err = executable.Exec(first); err != nil {
		t.Fatal(err)
	}
	if first.Buffer.String() != "first" {
		t.Fatalf("first output: %q", first.Buffer.String())
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	first.SetContext(canceled)
	first.AddError(errors.New("prior invocation"))
	pool.Put(first)
	// sync.Pool may discard this object, including deliberately under -race.
	// The contract is a clean, acquired state that can execute another request.
	next := pool.State()
	defer pool.Put(next)
	if next.Buffer.String() != "" || !next.IsValid() || next.Context().Err() != nil {
		t.Fatal("previous invocation state leaked")
	}
	if next.Take() {
		t.Fatal("returned state was not marked acquired")
	}
	if err = next.SetValue("foo", "second"); err != nil {
		t.Fatal(err)
	}
	if err = executable.Exec(next); err != nil {
		t.Fatal(err)
	}
	if next.Buffer.String() != "second" {
		t.Fatalf("second output: %q", next.Buffer.String())
	}
}
