package est

import "unsafe"

var trueValue = true
var falseValue = false
var TrueValuePtr = unsafe.Pointer(&trueValue)
var FalseValuePtr = unsafe.Pointer(&falseValue)

var emptyString = ""
var EmptyStringPtr = unsafe.Pointer(&emptyString)
