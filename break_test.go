package velty

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBreak_ForEach(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Items", []int{}))

	tpl := `#foreach($i in $Items)$i#if($i==2)#break#end#end`
	exec, newState, err := planner.Compile([]byte(tpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	_ = s.SetValue("Items", []int{1, 2, 3, 4})
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "12", s.Buffer.String())
}

func TestBreak_NearestLoopOnly(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Outer", []int{}))
	assert.Nil(t, planner.DefineVariable("Inner", []int{}))

	tpl := `#foreach($o in $Outer)#foreach($i in $Inner)$o:$i;#if($i==2)#break#end#end|#end`
	exec, newState, err := planner.Compile([]byte(tpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	_ = s.SetValue("Outer", []int{1, 2})
	_ = s.SetValue("Inner", []int{1, 2, 3})
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "1:1;1:2;|2:1;2:2;|", s.Buffer.String())
}

func TestBreak_OutsideLoop_Error(t *testing.T) {
	planner := New()
	_, _, err := planner.Compile([]byte(`#break`))
	assert.NotNil(t, err)
}
