package velty

import "testing"

type float32AssignmentRow struct {
	Value   float32
	Pointer *float32
	Next    *float32AssignmentRow
}

func TestFloat32IndirectAssignment(t *testing.T) {
	for _, value := range []float32{42.5, 0, -7.25} {
		for _, nilSource := range []bool{false, true} {
			previous := float32(99)
			root := &float32AssignmentRow{Value: value, Pointer: &value, Next: &float32AssignmentRow{Value: 99, Pointer: &previous}}
			if nilSource {
				root.Pointer = nil
			}
			planner := New()
			if err := planner.DefineVariable("Root", root); err != nil {
				t.Fatal(err)
			}
			executable, newState, err := planner.Compile([]byte("#set($Root.Next.Value = $Root.Value)#set($Root.Next.Pointer = $Root.Pointer)"))
			if err != nil {
				t.Fatal(err)
			}
			state := newState()
			if err = state.SetValue("Root", root); err != nil {
				t.Fatal(err)
			}
			if err = executable.Exec(state); err != nil {
				t.Fatal(err)
			}
			if root.Next.Value != value || root.Next.Pointer != root.Pointer {
				t.Fatalf("source=%v nil=%v: destination value=%v pointer=%p, wanted %p", value, nilSource, root.Next.Value, root.Next.Pointer, root.Pointer)
			}
		}
	}
}
