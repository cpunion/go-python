package gp

import "C"

type Char = C.char

//go:inline
func AllocCStr(s string) *Char {
	return C.CString(s)
}

func AllocCStrDontFree(s string) *Char {
	return C.CString(s)
}

func GoString(s *Char) string {
	return C.GoString((*Char)(s))
}

func GoStringN(s *Char, n int) string {
	return C.GoStringN(s, C.int(n))
}
