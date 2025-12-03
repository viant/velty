package velty

import (
	"github.com/stretchr/testify/assert"
	est "github.com/viant/velty/est"
	"testing"
)

// helpers
func compileAndRun(t *testing.T, tpl string, define func(p *Planner), set func(s *est.State)) (string, error) {
	// bring est into scope
	// note: alias import not allowed here; referencing fully-qualified in signature is enough
	p := New()
	if define != nil {
		define(p)
	}
	exec, newState, err := p.Compile([]byte(tpl))
	if err != nil {
		return "", err
	}
	st := newState()
	if set != nil {
		set(st)
	}
	if err := exec.Exec(st); err != nil {
		return "", err
	}
	return st.Buffer.String(), nil
}

// plain bool selector
func Test_If_Unary_Bool_Direct(t *testing.T) {
	tpl := `#if($B)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("B", false) },
		func(s *est.State) { _ = s.SetValue("B", true) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("B", false) },
		func(s *est.State) { _ = s.SetValue("B", false) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)
}

func Test_If_Unary_Bool_Negation(t *testing.T) {
	tpl := `#if(!$B)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("B", false) },
		func(s *est.State) { _ = s.SetValue("B", true) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("B", false) },
		func(s *est.State) { _ = s.SetValue("B", false) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)
}

// bool computed via function call
func Test_If_Unary_Bool_Computed(t *testing.T) {
	tpl := `#if($IsTrue("x"))T#elseF#end`
	p := New()
	_ = p.RegisterFunction("IsTrue", func(_ string) bool { return true })
	exec, newState, err := p.Compile([]byte(tpl))
	if !assert.Nil(t, err) {
		return
	}
	s := newState()
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "T", s.Buffer.String())

	// negation on computed bool
	tpl = `#if(!$IsTrue("x"))T#elseF#end`
	p = New()
	_ = p.RegisterFunction("IsTrue", func(_ string) bool { return false })
	exec, newState, err = p.Compile([]byte(tpl))
	if !assert.Nil(t, err) {
		return
	}
	s = newState()
	assert.Nil(t, exec.Exec(s))
	assert.Equal(t, "T", s.Buffer.String())
}

// string truthiness (non-empty true, empty false); negation not supported
func Test_If_Unary_String_Truthiness(t *testing.T) {
	tpl := `#if($S)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("S", "") },
		func(s *est.State) { _ = s.SetValue("S", "abc") },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("S", "") },
		func(s *est.State) { _ = s.SetValue("S", "") },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	// negation on non-bool should be a compile error
	_, err = compileAndRun(t, `#if(!$S)T#elseF#end`,
		func(p *Planner) { _ = p.DefineVariable("S", "") },
		func(s *est.State) { _ = s.SetValue("S", "abc") },
	)
	assert.NotNil(t, err)
}

// int truthiness (non-zero true, zero false); negation not supported
func Test_If_Unary_Int_Truthiness(t *testing.T) {
	tpl := `#if($N)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("N", 0) },
		func(s *est.State) { _ = s.SetValue("N", 3) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("N", 0) },
		func(s *est.State) { _ = s.SetValue("N", 0) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	_, err = compileAndRun(t, `#if(!$N)T#elseF#end`,
		func(p *Planner) { _ = p.DefineVariable("N", 0) },
		func(s *est.State) { _ = s.SetValue("N", 1) },
	)
	assert.NotNil(t, err)
}

// float truthiness (non-zero true, zero false)
func Test_If_Unary_Float_Truthiness(t *testing.T) {
	tpl := `#if($F)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("F", 0.0) },
		func(s *est.State) { _ = s.SetValue("F", 1.5) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("F", 0.0) },
		func(s *est.State) { _ = s.SetValue("F", 0.0) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	_, err = compileAndRun(t, `#if(!$F)T#elseF#end`,
		func(p *Planner) { _ = p.DefineVariable("F", 0.0) },
		func(s *est.State) { _ = s.SetValue("F", 2.0) },
	)
	assert.NotNil(t, err)
}

// slice truthiness (non-empty true, empty false)
func Test_If_Unary_Slice_Truthiness(t *testing.T) {
	tpl := `#if($Items)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Items", []int{}) },
		func(s *est.State) { _ = s.SetValue("Items", []int{1, 2}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Items", []int{}) },
		func(s *est.State) { _ = s.SetValue("Items", []int{}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	// nil slice should be false
	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Items", []int{}) },
		func(s *est.State) { _ = s.SetValue("Items", ([]int)(nil)) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	_, err = compileAndRun(t, `#if(!$Items)T#elseF#end`,
		func(p *Planner) { _ = p.DefineVariable("Items", []int{}) },
		func(s *est.State) { _ = s.SetValue("Items", []int{1}) },
	)
	assert.NotNil(t, err)
}

// pointer truthiness (non-nil true, nil false)
func Test_If_Unary_Pointer_Truthiness(t *testing.T) {
	type X struct{ V int }
	tpl := `#if($Ptr)T#elseF#end`
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*X)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", &X{V: 1}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)

	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*X)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", (*X)(nil)) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)
}

// zero-check via IsZero() bool on value receiver
type ziVal struct{ A int }

func (z ziVal) IsZero() bool { return z.A == 0 }

// zero-check via Zero() bool on pointer receiver
type ziPtr struct{ A int }

func (z *ziPtr) Zero() bool { return z.A == 0 }

func Test_If_Unary_PtrStruct_IsZero_ValueReceiver(t *testing.T) {
	tpl := `#if($Ptr)T#elseF#end`
	// nil pointer
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*ziVal)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", (*ziVal)(nil)) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	// non-nil ptr to zero value -> IsZero() == true => false
	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*ziVal)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", &ziVal{}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	// non-nil ptr to non-zero value -> IsZero() == false => true
	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*ziVal)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", &ziVal{A: 1}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)
}

func Test_If_Unary_PtrStruct_Zero_PtrReceiver(t *testing.T) {
	tpl := `#if($Ptr)T#elseF#end`
	// nil pointer
	out, err := compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*ziPtr)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", (*ziPtr)(nil)) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	// non-nil ptr to zero value -> Zero() == true => false
	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*ziPtr)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", &ziPtr{}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "F", out)

	// non-nil ptr to non-zero value -> Zero() == false => true
	out, err = compileAndRun(t, tpl,
		func(p *Planner) { _ = p.DefineVariable("Ptr", (*ziPtr)(nil)) },
		func(s *est.State) { _ = s.SetValue("Ptr", &ziPtr{A: 2}) },
	)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "T", out)
}

// unsupported comparable kinds: map, struct
func Test_If_Unary_Map_Struct_Unsupported(t *testing.T) {
	// map unsupported
	_, err := compileAndRun(t, `#if($M)T#elseF#end`,
		func(p *Planner) { _ = p.DefineVariable("M", map[string]int{}) },
		nil,
	)
	assert.NotNil(t, err)

	// struct unsupported
	type S struct{ A int }
	_, err = compileAndRun(t, `#if($Obj)T#elseF#end`,
		func(p *Planner) { _ = p.DefineVariable("Obj", S{}) },
		nil,
	)
	assert.NotNil(t, err)
}
