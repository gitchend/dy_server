package tools

import (
	"fmt"
	"unsafe"
)

func AssertNil(v interface{}) {
	if !IsNil(v) {
		panic(v)
	}
}

func AssertNotNil(v interface{}) {
	if IsNil(v) {
		panic(v)
	}
}

func AssertTrue(b bool, msg string, param ...interface{}) {
	if !b {
		info := fmt.Sprintf(msg, param...)
		panic(info)
	}
}
func IsNil(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
