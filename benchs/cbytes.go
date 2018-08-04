package benchs

import (
	"bytes"
	"strings"
)

var (
	prebytes  = []byte("hello")
	prestring = "hello"
)

// CompareBytes 2 slices of bytes
func CompareBytes(a, b []byte) bool {
	return bytes.Compare(a, b) == 0
}

// CompareStringByte between string and byte
func CompareStringByte(a []byte, b string) bool {
	return string(a) == b
}

// CompareStrings between 2 strings
func CompareStrings(a, b string) bool {
	return strings.Compare(a, b) == 0
}

// ComparePrecompiledByte with byte slice
func ComparePrecompiledByte(a []byte) bool {
	return bytes.Compare(a, prebytes) == 0
}

// ComparePrecompiledString with string
func ComparePrecompiledString(a string) bool {
	return strings.Compare(a, prestring) == 0
}

func CompareV1(data []byte) bool {
	return string(data) == "hello"
}

func CompareV2(data []byte) bool {
	return bytes.Compare(data, prebytes) == 0
}
