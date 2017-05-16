package format

import "unsafe"

func byteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
