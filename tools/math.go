package tools

import (
	"math/rand"
)

func ABS(num int32) int32 {
	y := num >> 31
	return (num ^ y) - y
}

func Int64Merge(h, l int32) (id int64) {
	return int64(uint64(h)<<32 | uint64(l))
}

func Int64Split(id int64) (h, l int32) {
	return int32(uint64(id) >> 32), int32(uint32(id))
}

func Int32Merge(h, l int16) (id int32) {
	return int32(uint32(h)<<16 | uint32(l))
}

func Int32Split(id int32) (h, l int16) {
	return int16(uint32(id) >> 16), int16(uint16(id))
}

func UInt16Merge(h, l uint8) (id uint16) {
	return uint16(h)<<8 | uint16(l)
}

func UInt16Split(id uint16) (h, l uint8) {
	return uint8(id >> 8), uint8(id)
}

//闭区间, 结果不重复
func RandBetween(start int, end int, count int) []int {
	existMap := make(map[int]bool)

	if (end-start)+1 < count {
		count = end - start + 1
	}

	if end-start == 0 {
		return []int{start}
	}

	nums := make([]int, 0)

	for len(nums) < count {
		num := rand.Intn((end-start)+1) + start
		if _, exist := existMap[num]; !exist {
			nums = append(nums, num)
			existMap[num] = true
		}
	}

	return nums
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt32(v int32, others ...int32) int32 {
	for _, other := range others {
		if v < other {
			v = other
		}
	}
	return v
}

func MinInt32(v int32, others ...int32) int32 {
	for _, other := range others {
		if v > other {
			v = other
		}
	}
	return v
}

func MaxInt64(v int64, others ...int64) int64 {
	for _, other := range others {
		if v < other {
			v = other
		}
	}
	return v
}

func MinInt64(v int64, others ...int64) int64 {
	for _, other := range others {
		if v > other {
			v = other
		}
	}
	return v
}
