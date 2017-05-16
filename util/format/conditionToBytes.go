package format

import (
	"reflect"
	"unsafe"
)

type Condition struct {
	EnrollNum   int
	AvgScore    int
	LowScore    int
	Rank        int
	Major       string
	EnrollScore []int
	IsSuccess   bool
}

var sizeOfCondition = int(unsafe.Sizeof(Condition{}))

func ConditionToBytes(s *[]Condition) []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfCondition
	x.Cap = sizeOfCondition
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}
