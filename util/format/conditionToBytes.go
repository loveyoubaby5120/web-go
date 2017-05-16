package format

import (
	"reflect"
	"unsafe"
)

// Condition is hight school score config
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

// ConditionToBytes is format Condition to bytes
func ConditionToBytes(s *[]Condition) []byte {
	var x reflect.SliceHeader
	x.Len = sizeOfCondition
	x.Cap = sizeOfCondition
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}
