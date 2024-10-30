package gp

/*
#include <Python.h>
*/
import "C"
import (
	"unsafe"
)

//go:inline
func AllocCStr(s string) *C.char {
	return C.CString(s)
}

func AllocWCStr(s string) *C.wchar_t {
	runes := []rune(s)
	wchars := make([]uint16, len(runes)+1)
	for i, r := range runes {
		wchars[i] = uint16(r)
	}
	wchars[len(runes)] = 0
	return (*C.wchar_t)(unsafe.Pointer(&wchars[0]))
}

func GoString(s *C.char) string {
	return C.GoString((*C.char)(s))
}

func GoStringN(s *C.char, n int) string {
	return C.GoStringN((*C.char)(s), C.int(n))
}
