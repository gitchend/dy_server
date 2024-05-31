package tools

import (
	"fmt"
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
