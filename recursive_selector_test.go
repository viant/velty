package velty

import (
	"fmt"
	"github.com/viant/velty/ast/expr"
	"reflect"
	"strings"
	"sync"
	"testing"
)

type selectorNode struct {
	Allowed  bool
	Accepted *bool
	Padding  string
	Value    int `velty:"names=Amount"`
	Next     *selectorNode
	Children []*selectorNode
	Values   []selectorNode
	Peers    map[string]*selectorNode
	Hidden   int `velty:"-"`
}

func (n *selectorNode) Number() int {
	if n == nil {
		return -1
	}
	return n.Value
}
func (n *selectorNode) Child() *selectorNode {
	if n == nil {
		return nil
	}
	return n.Next
}
func (n *selectorNode) Add(value int) int { return n.Value + value }

type selectorLeft struct {
	Name  string
	Right *selectorRight
}
type selectorRight struct {
	Value int
	Left  *selectorLeft
}
type selectorSiblings struct {
	Padding int
	Left    *selectorNode
	Right   *selectorNode
}

func TestRecursiveSelectors(t *testing.T) {
	tail := &selectorNode{Value: 37}
	middle := &selectorNode{Value: 23, Next: tail, Children: []*selectorNode{tail}, Values: []selectorNode{*tail}, Peers: map[string]*selectorNode{"x": tail}}
	root := &selectorNode{Value: 11, Next: middle, Children: []*selectorNode{middle, tail}, Values: []selectorNode{*middle}, Peers: map[string]*selectorNode{"x": middle}}
	mutual := &selectorLeft{Name: "left"}
	mutual.Right = &selectorRight{Value: 41, Left: mutual}
	for _, tc := range []struct {
		name, template, want string
		input                interface{}
	}{
		{"pointer", "$Root.Next.Next.Amount", "37", root},
		{"method", "$Root.Next.Next.Number()", "37", root},
		{"method argument", "$Root.Next.Add($Root.Next.Next.Amount)", "60", root},
		{"method result", "$Root.Child().Child().Amount", "37", root},
		{"method after result field", "$Root.Child().Next.Number()", "37", root},
		{"slice", "$Root.Children[0].Children[0].Amount", "37", root},
		{"value slice", "$Root.Values[0].Values[0].Amount", "37", root},
		{"slice method", "$Root.Children[1].Number()", "37", root},
		{"slice pointer method", "$Root.Children[0].Next.Number()", "37", root},
		{"map method", `${Root.Peers["x"].Number()}`, "23", root},
		{"nil indexed method", "$Root.Children[0].Number()", "-1", &selectorNode{Children: []*selectorNode{nil}}},
		{"nil returned method", "$Root.Next.Next.Child().Number()", "-1", root},
		{"map", `${Root.Peers["x"].Peers["x"].Amount}`, "37", root},
		{"foreach", "#foreach($item in $Root.Children)$item.Next.Amount;#end", "37;0;", root},
		{"nil field", "$Root.Next.Next.Next.Amount", "0", root},
		{"nil root", "$Root.Next.Next.Amount", "0", (*selectorNode)(nil)},
		{"nil method", "$Root.Next.Next.Next.Number()", "-1", root},
		{"nil method result", "$Root.Next.Next.Child().Amount", "0", root},
		{"nil truth", "#if($Root.Next.Next.Next)yes#else no#end", " no", root},
		{"mutual", "$Root.Right.Left.Right.Left.Right.Value", "41", mutual},
		{"siblings", "$Root.Left.Next.Amount/$Root.Right.Next.Amount", "23/37", &selectorSiblings{Left: root, Right: middle}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := New()
			if err := p.DefineVariable("Root", tc.input); err != nil {
				t.Fatal(err)
			}
			executable, newState, err := p.Compile([]byte(tc.template))
			if err != nil {
				t.Fatal(err)
			}
			state := newState()
			if err = state.SetValue("Root", tc.input); err != nil {
				t.Fatal(err)
			}
			if err = executable.Exec(state); err != nil {
				t.Fatal(err)
			}
			if got := state.Buffer.String(); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestRecursiveSelectorsAuthoredDepth(t *testing.T) {
	for _, depth := range []int{1, 32, 257} {
		t.Run(fmt.Sprint(depth), func(t *testing.T) {
			root := &selectorNode{Value: 73}
			root.Next = root
			p := New()
			if err := p.DefineVariable("Root", root); err != nil {
				t.Fatal(err)
			}
			registered := len(p.selectors.Index)
			executable, newState, err := p.Compile([]byte("$Root" + strings.Repeat(".Next", depth) + ".Amount"))
			if err != nil {
				t.Fatal(err)
			}
			if len(p.selectors.Index) != registered {
				t.Fatal("authored recursive paths expanded the global registry")
			}
			var wg sync.WaitGroup
			for i := 0; i < 8; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					state := newState()
					if err := state.SetValue("Root", root); err != nil {
						t.Error(err)
						return
					}
					if err := executable.Exec(state); err != nil {
						t.Error(err)
						return
					}
					if state.Buffer.String() != "73" {
						t.Error(state.Buffer.String())
					}
				}()
			}
			wg.Wait()
		})
	}
}

func TestRecursiveSelectorsInvalidField(t *testing.T) {
	for _, path := range []string{"$Root.Next.Next.Missing", "$Root.Next.Next.Hidden", "$Root.Next.Next.Amount.Missing", "$Root.Next.Next.Value"} {
		p := New()
		if err := p.DefineVariable("Root", &selectorNode{}); err != nil {
			t.Fatal(err)
		}
		if _, _, err := p.Compile([]byte(path)); err == nil {
			t.Fatalf("accepted %s", path)
		}
	}
}

type recursiveFactory struct{ Root *selectorNode }

func (f *recursiveFactory) Give() interface{} { return f.Root }
func (f *recursiveFactory) MethodResultType(_ string, _ *expr.Call) (reflect.Type, error) {
	return reflect.TypeOf(f.Root), nil
}
func TestRecursiveInferredMethodResult(t *testing.T) {
	for _, root := range []*selectorNode{{Value: 17, Next: &selectorNode{Value: 29}}, nil} {
		p := New()
		if err := p.RegisterFuncNs("factory", &recursiveFactory{Root: root}); err != nil {
			t.Fatal(err)
		}
		executable, newState, err := p.Compile([]byte("#set($local = $factory.Give())$local.Next.Amount|$factory.Give().Next.Amount"))
		if err != nil {
			t.Fatal(err)
		}
		state := newState()
		if err := executable.Exec(state); err != nil {
			t.Fatal(err)
		}
		want := "29|29"
		if root == nil {
			want = "0|0"
		}
		if state.Buffer.String() != want {
			t.Fatalf("got %q want %q", state.Buffer.String(), want)
		}
	}
}

type recursiveEmbedded struct {
	Padding       int
	*selectorNode `velty:"prefix=Node"`
	Next          *recursiveEmbedded
}

func TestRecursiveEmbeddedSelectors(t *testing.T) {
	root := &recursiveEmbedded{selectorNode: &selectorNode{Value: 19}}
	root.Next = root
	root.selectorNode.Next = root.selectorNode
	for _, template := range []string{"$Root.Next.Next.NodeAmount", "$Root.Next.NodeNext.Amount"} {
		p := New()
		if err := p.DefineVariable("Root", root); err != nil {
			t.Fatal(err)
		}

		executable, newState, err := p.Compile([]byte(template))
		if err != nil {
			t.Fatal(err)
		}
		state := newState()
		if err := state.SetValue("Root", root); err != nil {
			t.Fatal(err)
		}
		if err := executable.Exec(state); err != nil {
			t.Fatal(err)
		}
		if state.Buffer.String() != "19" {
			t.Fatal(state.Buffer.String())
		}
	}
}

func TestCycleDetectorAncestry(t *testing.T) {
	root := NewCycleDetector(reflect.TypeOf(selectorSiblings{}))
	left, cycle := root.Child(reflect.TypeOf(&selectorNode{}), nil)
	if cycle {
		t.Fatal("first child marked cyclic")
	}
	right, cycle := root.Child(reflect.TypeOf([]selectorNode{}), nil)
	if cycle || left == right {
		t.Fatal("independent sibling reused ancestry")
	}
	if _, cycle := left.Child(reflect.TypeOf(selectorNode{}), nil); !cycle {
		t.Fatal("pointer recursion not detected")
	}
	if _, cycle := right.Child(reflect.TypeOf(map[string]*selectorNode{}), nil); !cycle {
		t.Fatal("container recursion not detected")
	}
}

type recursiveList []recursiveList

func (l recursiveList) Size() int { return len(l) }

func TestRecursiveContainerType(t *testing.T) {
	root := recursiveList{recursiveList{nil, nil}}
	p := New()
	if err := p.DefineVariable("Root", root); err != nil {
		t.Fatal(err)
	}
	executable, newState, err := p.Compile([]byte("$Root[0].Size()"))
	if err != nil {
		t.Fatal(err)
	}
	state := newState()
	if err := state.SetValue("Root", root); err != nil {
		t.Fatal(err)
	}
	if err := executable.Exec(state); err != nil {
		t.Fatal(err)
	}
	if state.Buffer.String() != "2" {
		t.Fatal(state.Buffer.String())
	}
}

type recursiveContainerFactory struct{ Value interface{} }

func (f *recursiveContainerFactory) Give() interface{} { return f.Value }
func (f *recursiveContainerFactory) MethodResultType(_ string, _ *expr.Call) (reflect.Type, error) {
	return reflect.TypeOf(f.Value), nil
}
func TestRecursiveInferredContainerResult(t *testing.T) {
	node := &selectorNode{Next: &selectorNode{Value: 43}}
	for _, tc := range []struct {
		value interface{}
		path  string
	}{
		{map[string]*selectorNode{"key": node}, `["key"].Next.Amount`},
		{[]*selectorNode{node}, `[0].Next.Amount`},
		{*node, `.Next.Amount`},
	} {
		p := New()
		if err := p.RegisterFuncNs("factory", &recursiveContainerFactory{Value: tc.value}); err != nil {
			t.Fatal(err)
		}
		executable, newState, err := p.Compile([]byte("${factory.Give()" + tc.path + "}"))
		if err != nil {
			t.Fatal(err)
		}
		state := newState()
		if err := executable.Exec(state); err != nil {
			t.Fatal(err)
		}
		if state.Buffer.String() != "43" {
			t.Fatalf("%T: %q", tc.value, state.Buffer.String())
		}
	}
}

func TestRecursiveAssignmentAndReuse(t *testing.T) {
	p := New()
	if err := p.DefineVariable("Root", &selectorNode{}, "Alias"); err != nil {
		t.Fatal(err)
	}
	executable, newState, err := p.Compile([]byte("#set($Root.Next.Next.Amount = 91)$Alias.Next.Next.Amount"))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		root := &selectorNode{Next: &selectorNode{Next: &selectorNode{Value: i}}}
		state := newState()
		if err := state.SetValue("Root", root); err != nil {
			t.Fatal(err)
		}
		if err := executable.Exec(state); err != nil {
			t.Fatal(err)
		}
		if state.Buffer.String() != "91" || root.Next.Next.Value != 91 {
			t.Fatalf("assignment/reuse: %q %+v", state.Buffer.String(), root.Next.Next)
		}
	}
}

func TestRecursiveBooleanAssignment(t *testing.T) {
	yes, no := true, false
	root := &selectorNode{Allowed: true, Accepted: &yes, Next: &selectorNode{Accepted: &no}}
	p := New()
	if err := p.DefineVariable("Root", root); err != nil {
		t.Fatal(err)
	}
	executable, newState, err := p.Compile([]byte("#set($Root.Next.Allowed = $Root.Allowed)#set($Root.Next.Accepted = $Root.Accepted)#if($Root.Next.Allowed)yes#end"))
	if err != nil {
		t.Fatal(err)
	}
	state := newState()
	if err := state.SetValue("Root", root); err != nil {
		t.Fatal(err)
	}
	if err := executable.Exec(state); err != nil {
		t.Fatal(err)
	}
	if !root.Next.Allowed || root.Next.Accepted != &yes || state.Buffer.String() != "yes" {
		t.Fatalf("bool assignment: %+v %q", root.Next, state.Buffer.String())
	}
}
