# velty (template evaluator in go)

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

## Introduction

In order to reduce execution time, this project first produces execution plan alongside with all variables needed to execute it.
One execution plan can be shared alongside many instances of scoped variables needed by executor.
Scoped Variables holds both execution state  and  variables defined or used in the evaluation code.

```go
    planner := plan.New(8192)
    executor, newState, err := planner.Compile(code)
    vars := newState() //creates memory instance needed by executor 
    executor.Exec(vars)
```

## Usage

### Expression

## Performance

### Expression evaluation

### Code execution

## Bugs

## Contributing to igo

Velty is an open source project and contributors are welcome!

See [TODO](TODO.md) list

## Credits and Acknowledgements
