package format

import "unsafe"

// ByteToString is format byte to String
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
