package tools

import (
	"reflect"
	"unsafe"
)

//重置数据
func Reset(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

//判空
func IsNil(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
