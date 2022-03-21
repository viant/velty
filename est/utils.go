package est

import "unsafe"

var trueValue = true
var falseValue = false
var TrueValuePtr = unsafe.Pointer(&trueValue)
var FalseValuePtr = unsafe.Pointer(&falseValue)

var emptyString = ""
var EmptyStringPtr = unsafe.Pointer(&emptyString)

var zeroInt = 0
var ZeroIntPtr = unsafe.Pointer(&zeroInt)

var zeroFloat = 0.0
var ZeroFloatPtr = unsafe.Pointer(&zeroFloat)
