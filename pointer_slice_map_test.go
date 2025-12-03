package velty

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type subPS struct{ Title string }
type itemPS struct{ Sub *subPS }
type modelPS struct{ Items []*itemPS }
type modelNums struct{ Nums []*int }
type itemMap struct{ Sub *subPS }
type modelMap struct{ M map[string]*itemMap }

func Test_SliceOfPointer_ElementField(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Model", modelPS{}))
	tmpl := `$Model.Items[0].Sub.Title`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	hello := &subPS{Title: "Hello"}
	s.SetValue("Model", modelPS{Items: []*itemPS{{Sub: hello}}})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "Hello", s.Buffer.String())
}

func Test_SliceOfPointer_PrimitivePointer(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Model", modelNums{}))

	tmpl := `$Model.Nums[0]`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}

	v := 5
	s := newState()
	s.SetValue("Model", modelNums{Nums: []*int{&v}})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "5", s.Buffer.String())
}

func Test_MapOfPointer_ElementField(t *testing.T) {
	planner := New()
	assert.Nil(t, planner.DefineVariable("Model", modelMap{}))

	tmpl := `${Model.M["k"].Sub.Title}`
	exec, newState, err := planner.Compile([]byte(tmpl))
	if !assert.Nil(t, err) {
		return
	}

	s := newState()
	s.SetValue("Model", modelMap{M: map[string]*itemMap{"k": {Sub: &subPS{Title: "X"}}}})
	err = exec.Exec(s)
	assert.Nil(t, err)
	assert.Equal(t, "X", s.Buffer.String())
}
