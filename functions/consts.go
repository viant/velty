package functions

import "reflect"

var FuncErrors = registryInstance.DefineNs("errors", NewEntry(
	&Errors{},
	NewFunctionNamespace(reflect.TypeOf(&Errors{})),
))

var FuncTime = registryInstance.DefineNs("time", NewEntry(
	&Time{},
	NewFunctionNamespace(reflect.TypeOf(&Time{})),
))

var FuncTypes = registryInstance.DefineNs("types", NewEntry(
	&Types{},
	NewFunctionNamespace(reflect.TypeOf(&Types{})),
))

var FuncSlices = registryInstance.DefineNs("slices", NewEntry(
	&Slices{},
	NewFunctionNamespace(reflect.TypeOf(&Slices{})),
))

var FuncStrconv = registryInstance.DefineNs("strconv", NewEntry(
	&Strconv{},
	NewFunctionNamespace(reflect.TypeOf(&Strconv{})),
))

var FuncMath = registryInstance.DefineNs("math", NewEntry(
	&Math{},
	NewFunctionNamespace(reflect.TypeOf(&Math{})),
))

var FuncStrings = registryInstance.DefineNs("strings", NewEntry(
	&Strings{},
	NewFunctionNamespace(reflect.TypeOf(&Strings{})),
))

var FuncMaps = registryInstance.DefineNs("maps", NewEntry(
	&Maps{},
	NewFunctionNamespace(reflect.TypeOf(&Maps{})),
))

var MapHasKey = registryInstance.DefineNs("HasKey", NewEntry(
	HasKeyFunc.handler,
	NewFunctionKind([]reflect.Kind{HasKeyFunc.kind}),
))

var SliceIndexBy = registryInstance.DefineNs("IndexBy", NewEntry(
	SliceIndexByFunc.Handler(),
	NewFunctionKind(SliceIndexByFunc.Kind()),
))
