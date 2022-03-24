# Velty (template engine in go)

[![GoReportCard](https://goreportcard.com/badge/github.com/viant/velty)](https://goreportcard.com/report/github.com/viant/velty)
[![GoDoc](https://godoc.org/github.com/viant/velty?status.svg)](https://godoc.org/github.com/viant/velty)

This library is compatible with Go 1.17+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Motivation](#motivation)
- [Introduction](#introduction)
- [Usage](#usage)
- [Performance](#performance)
- [Bugs](#bugs)
- [Contribution](#contributing-to-igo)
- [License](#license)

## Motivation

The goal of this library is to create and fill templates with data, using reflection and allocate new data as little as
possible. It implements subset of the java velocity library.

[//]: # (- [GoEval]&#40;https://github.com/xtaci/goeval&#41;)

[//]: # (- [GoVal]&#40;https://github.com/maja42/goval&#41;)

[//]: # (- [Yaegi]&#40;https://github.com/traefik/yaegi&#41; .)

See [performance](#performance) section for details.

## Introduction

In order to reduce execution time, this project first produces execution plan alongside with all variables needed to
execute it. One execution plan can be shared alongside many instances of scoped variables needed by executor. Scoped
Variables holds both execution state and variables defined or used in the evaluation code.

```go
    planner := plan.New()
    exec, newState, err := planner.Compile(code)
   
    state := newState() 
    exec.Exec(state)
    fmt.Printf("Result: %v", state.Buffer.String())
   
    anotherState := newState()
    exec.Exec(anotherState)
    fmt.Printf("Result: %v", anotherState.Buffer.String())
```

## Usage

In order to create execution plan, you need to create a planner:
```go
    planner := plan.New()
```

You can also specify the initial buffer and cache size. If you want to do that, you can pass optional arguments in given
order: `BufferSize` -> `CacheSize`.

```go
    planner := plan.New(1024, 200)
```

Once you have the `Planner` you have to define variables that will be used. Velty doesn't use a `map` to store state, but it 
recreates an internal type each time you define new variable and uses `reflect.StructField.Offset` to access data from the state.
Velty supports two ways of defining planner variables:

* `planner.DefineVariable(variableName, variableType)` - will create and add non-anonymous `reflect.StructField`
* `planner.EmbedVariable(variableName, variableType)` - will create and add anonymous `reflect.StructField`

For each of the non-anonymous struct field registered with `DefineVariable` or `EmbedVariable` will be created unique `Selector`.
Selector is used to get field value from the state. 

```go
  err = planner.DefineVariable("foo", reflect.Typeof(Foo{})) 
  //handle error if needed
  err = planner.DefineVariable("boo", Boo{}) 
  //handle error if needed
  err = planner.EmbedVariable("bar", reflect.Typeof(Bar{})) 
  //handle error if needed
  err = planner.EmbedVariable("emp", Boo{}) 
  //handle error if needed
```

You can pass an instance or the `reflect.Type`. However, there are some constraints:
* Velty creates selector for each of the struct field. If you define i.e.:
```go
  type Foo struct {
      Name string
      ID int
  }
  
  planner.DefineVariable("foo", Foo{})
```
Velty will create three selectors: `foo`, `foo___Name`, `foo_ID`. Structure used by the velty shouldn't have three consecutive 
underscores in any of the fields.

* Velty won't create selectors for the Anonymous fields and will flatten the fields of the anonymous field.
```go
  type Foo struct {
      Name string
      ID int
  }
  
  type Bar struct {
      Foo
  }

  planner.EmbedVariable("foo", Bar{})
```
Velty will create only two selectors: `Name` and `ID` because all other fields are Anonymous.

* You can use tags to customize selector id, see: [Tags](#tags)
* Velty generates selectors for the constants and name them: `_T0`, `_T1`, `_T2` etc.

In the next step you can register functions. In the template you use the receiver syntax
i.e. `foo.Name.ToUpperCase()` but in the Go, you have to register plain function, where the first argument is the value of
field on which function was called.

```go
  err = planner.RegisterFunction("ToUpperCase", strings.ToUpper) 
  //handle error if needed
```

You can register function in two ways:
* `planner.RegisterFunction` - you can register regular functions like `strings.ToUpper`, and some of them are optimized using
type assertion. If the function isn't optimized, it will be called via `reflect.ValueOf.Call`. 

* `planner.RegsiterFunc` - if you notice that function is not optimized, you can optimize it registering `*op.Func`. 
The simple implementation:

```go
    customFunc := &op.Func{
		ResultType: reflect.TypeOf(""),
		Function: func(accumulator *Selector, operands []*Operand, state *est.State) unsafe.Pointer {
			if len(operands) < 2 {
				return nil
			}
			
			accumulator.SetBool(state.MemPtr, strings.HasPrefix(*(*string)(operands[0].Exec(state)), *(*string)(operands[1].Exec(state))))
			return accumulator.Pointer(state.MemPtr)
		},
	}
	
    err = planner.RegisterFunc("HasPrefix", customFunc) 
    //handle error if needed
```

Regular function can return no more than two non-pointer values. First is the new value, the second is an error. 
However errors in this case are ignored, and if any returned - the zero value will be appended to the result. 

The next step is to create execution plan and new state function:
```go
  template := `...`
  exec, newState, err := planner.Compile([]byte(template)) 
  // handle error if needed
  state := newState()
  exec.Exec(state)
```

### Tags
In order to match template identifiers with the struct fields, you can use the `velty` tag. 
Supported attributes:
* `name` - represents template identifier name i.e.:
```go
  type Foo struct {
    Name string `velty:"name=fooName"`
  }

  planner.DefineVariable("foo", Foo{})
  template := `${foo.fooName}`
```
* `names` - similar to the `name` but in this case you can specify more than one template identifier by separating them with `|`
```go
  type Foo struct {
    Name string `velty:"name=NAME|FOO_NAME"`
  }
   
  planner.DefineVariable("foo", Foo{})
  template := `${foo.NAME}, ${foo.FOO_NAME}`
```
* `prefix` - prefix can be used on the anonymous fields:
```go
  type Foo struct {
      Name string `velty:"name=NAME"`
  }
    
  type Boo struct {
      Foo `velty:"prefix=FOO_"`
  }

  planner.EmbedVariable("boo", Boo{})
  template := `${FOO_NAME}`
```

* `-` - tells Velty to don't create selector:
```go
 type Foo struct {
      Name string `velty:"-"`
  }

  planner.EmbedVariable("foo", Foo{})
  template := `${foo.Name}` // throws an error during compile time
```
## Bugs

This project does not implement full java velocity spec, but just a subset. It supports:
* variables - i.e. `${foo.Name} $Name`
* assignment - i.e. `#set($var1 = 10 + 20 * 10) #set($var2 = ${foo.Name})`
* if statements - i.e. `#if(1==1) abc #elsif(2==2) def #else ghi #end`
* foreach - i.e. `#foreach($name in ${foo.Names})`
* function calls - i.e. `${name.toUpper()}`
* template evaluation - i.e. `#evaluate($TEMPLATE)`

## Contributing to Velty

Velty is an open source project and contributors are welcome!
