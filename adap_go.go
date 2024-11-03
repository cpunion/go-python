package gp

/*
#include <Python.h>
*/
import "C"
import (
	"unsafe"
)

type Char = C.char
type WChar = C.wchar_t
type Int = C.int
type Pointer = unsafe.Pointer

//go:inline
func AllocCStr(s string) *C.char {
	return C.CString(s)
}

func AllocCStrDontFree(s string) *C.char {
	return C.CString(s)
}

func GoString(s *C.char) string {
	return C.GoString((*C.char)(s))
}

func GoStringN(s *C.char, n int) string {
	return C.GoStringN((*C.char)(s), C.int(n))
}
