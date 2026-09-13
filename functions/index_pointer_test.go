package functions

import (
	"reflect"
	"testing"
)

func TestIndexByNormalizesPointerKeys(t *testing.T) {
	type row struct {
		ID   *int
		Name string
	}
	first, second := 7, 7
	for _, values := range []interface{}{[]row{{ID: &first, Name: "present"}}, []*row{{ID: &first, Name: "present"}}} {
		index, err := SliceIndexByFunc.Handler().(func(interface{}, string) (interface{}, error))(values, "ID")
		if err != nil {
			t.Fatal(err)
		}
		if value := reflect.ValueOf(index).MapIndex(reflect.ValueOf(7)); !value.IsValid() {
			t.Fatalf("index did not normalize pointer key for %T: %#v", values, index)
		}
		found, err := HasKeyFunc.handler.(func(interface{}, interface{}) (bool, error))(index, &second)
		if err != nil || !found {
			t.Fatalf("equivalent key pointer lookup: %v %v", found, err)
		}
		var absent *int
		found, err = HasKeyFunc.handler.(func(interface{}, interface{}) (bool, error))(index, absent)
		if err != nil || found {
			t.Fatalf("nil key lookup %v %v", found, err)
		}
	}
}

func TestIndexByPreservesUnsignedPointerKeyIdentity(t *testing.T) {
	type row struct{ ID *uint64 }
	first, second, other := ^uint64(0), ^uint64(0), ^uint64(0)-1
	index, err := SliceIndexByFunc.Handler().(func(interface{}, string) (interface{}, error))([]*row{{ID: &first}}, "ID")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		key  interface{}
		want bool
	}{{&second, true}, {second, true}, {&other, false}, {(*uint64)(nil), false}} {
		found, err := HasKeyFunc.handler.(func(interface{}, interface{}) (bool, error))(index, tt.key)
		if err != nil || found != tt.want {
			t.Fatalf("key %v found %v err %v", tt.key, found, err)
		}
	}
}
