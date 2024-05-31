package tools

import (
	"reflect"
)

func UniqueSlice(slice *[]interface{}) {
	found := make(map[interface{}]bool)
	total := 0
	for i, val := range *slice {
		if _, ok := found[val]; !ok {
			found[val] = true
			(*slice)[total] = (*slice)[i]
			total++
		}
	}

	*slice = (*slice)[:total]
}

func InsertToSlice(slice, insertion []interface{}, index int) []interface{} {
	result := make([]interface{}, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}

func AppendToSlice(slice, appendItems []interface{}) []interface{} {
	return InsertToSlice(slice, appendItems, len(slice))
}

func ReverseSlice(slice []interface{}) []interface{} {
	if len(slice) == 0 || len(slice) == 1 {
		return slice
	}
	return append(ReverseSlice(slice[1:]), slice[0])
}

func RemoveFromSlice(slice []interface{}, start int, end int) []interface{} {
	return append(slice[:start], slice[end+1:]...)
}

// condition is the function which is used to check if the element shall be removed
func RemoveFromSliceIf(slice []interface{}, condition func(interface{}) bool) []interface{} {
	result := make([]interface{}, len(slice))
	for _, item := range slice {
		// if matched, just ignore
		if !condition(item) {
			result = append(result, item)
		}
	}
	return result
}

func Int8SliceSearch(s []int8, t int8) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}

	return false
}

func Int32SliceSearch(s []int32, t int32) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}

	return false
}

func Int64SliceSearch(s []int64, t int64) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}

	return false
}

func StringSliceSearch(s []string, t string) bool {
	for _, n := range s {
		if n == t {
			return true
		}
	}

	return false
}

func Reverse(slice interface{}, size func() int, swap func(i, j int)) {
	if swap == nil || size == nil {
		return
	}

	for i := 0; i < size()/2; i++ {
		swap(i, size()-i-1)
	}
}

func CheckDuplicate(slice interface{}, check func(i, j int) bool) bool {
	if slice == nil || check == nil {
		return true
	}

	rv := reflect.ValueOf(slice)
	len := rv.Len()
	for i := 0; i < len; i++ {
		for j := i + 1; j < len; j++ {
			if check(i, j) {
				return true
			}
		}
	}
	return false
}
