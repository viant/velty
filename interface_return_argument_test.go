package velty

import (
	"fmt"
	"testing"
)

type interfaceResultRecord struct {
	ID   int
	Name string
}
type interfaceResultMethods struct{}

func (*interfaceResultMethods) Give(index int) (interface{}, error) {
	switch index {
	case 0:
		return nil, nil
	case 1:
		return &interfaceResultRecord{ID: 7, Name: "record"}, nil
	case 2:
		return map[string]int{"id": 9}, nil
	}
	return (*interfaceResultRecord)(nil), nil
}
func (*interfaceResultMethods) Take(value interface{}) (string, error) {
	switch actual := value.(type) {
	case nil:
		return "nil", nil
	case *interfaceResultRecord:
		if actual == nil {
			return "typed-nil", nil
		}
		return fmt.Sprint(actual.ID), nil
	case map[string]int:
		return fmt.Sprint(actual["id"]), nil
	}
	return "", fmt.Errorf("unexpected %T", value)
}

func TestNestedInterfaceReturnArgument(t *testing.T) {
	planner := New()
	if err := planner.DefineVariable("Methods", &interfaceResultMethods{}); err != nil {
		t.Fatal(err)
	}
	execution, newState, err := planner.Compile([]byte(`$Methods.Take($Methods.Give(1))|$Methods.Take($Methods.Give(0))|$Methods.Take($Methods.Give(2))|$Methods.Take($Methods.Give(3))`))
	if err != nil {
		t.Fatal(err)
	}
	state := newState()
	if err = state.SetValue("Methods", &interfaceResultMethods{}); err != nil {
		t.Fatal(err)
	}
	if err = execution.Exec(state); err != nil {
		t.Fatal(err)
	}
	if actual := state.Buffer.String(); actual != "7|nil|9|typed-nil" {
		t.Fatalf("nested result %q", actual)
	}
}
