package velty

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fcItem struct{ Name string }

func Test_ForEach_Context_WithSliceOfStructs(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Items", []fcItem{}))

	tmpl := `#foreach($i in $Items)$i.Name#if($foreach.HasNext),#end#end`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	_ = s.SetValue("Items", []fcItem{{"A"}, {"B"}, {"C"}})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "A,B,C", s.Buffer.String())

	// Assert all foreach fields
	tmpl2 := `#foreach($i in $Items)[$foreach.Index,$foreach.Count,$foreach.First,$foreach.Last,$foreach.HasNext]#end`
	exec2, newState2, err := planner.Compile([]byte(tmpl2))
	if !assert.Nil(t, err) {
		return
	}
	s2 := newState2()
	_ = s2.SetValue("Items", []fcItem{{"A"}, {"B"}, {"C"}})
	err = exec2.Exec(s2)
	assert.Nil(t, err)
	assert.Equal(t, "[0,1,true,false,true][1,2,false,false,true][2,3,false,true,false]", s2.Buffer.String())
}

func Test_ForEach_Context_WithSliceOfPointers(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Items", []*fcItem{}))

	tmpl := `#foreach($i in $Items)$i.Name#if($foreach.HasNext),#end#end`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	a, b, c := &fcItem{"A"}, &fcItem{"B"}, &fcItem{"C"}
	_ = s.SetValue("Items", []*fcItem{a, b, c})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "A,B,C", s.Buffer.String())

	// Assert all foreach fields
	tmpl2 := `#foreach($i in $Items)[$foreach.Index,$foreach.Count,$foreach.First,$foreach.Last,$foreach.HasNext]#end`
	exec2, newState2, err := planner.Compile([]byte(tmpl2))
	if !assert.Nil(t, err) {
		return
	}
	s2 := newState2()
	a2, b2, c2 := &fcItem{"A"}, &fcItem{"B"}, &fcItem{"C"}
	_ = s2.SetValue("Items", []*fcItem{a2, b2, c2})
	err = exec2.Exec(s2)
	assert.Nil(t, err)
	assert.Equal(t, "[0,1,true,false,true][1,2,false,false,true][2,3,false,true,false]", s2.Buffer.String())
}

func Test_ForEach_Context_EmptySlice(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Items", []fcItem{}))

	// Even though we reference $foreach fields inside the loop, with an empty slice the body never runs.
	tmpl := `#foreach($i in $Items)[$foreach.Index,$foreach.Count,$foreach.First,$foreach.Last,$foreach.HasNext]#end`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	_ = s.SetValue("Items", []fcItem{})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "", s.Buffer.String())
}

func Test_ForEach_Context_EmptySlice_Pointers(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Items", []*fcItem{}))

	tmpl := `#foreach($i in $Items)[$foreach.Index,$foreach.Count,$foreach.First,$foreach.Last,$foreach.HasNext]#end`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	_ = s.SetValue("Items", []*fcItem{})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "", s.Buffer.String())
}
