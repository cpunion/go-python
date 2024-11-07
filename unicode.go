package gp

/*
#include <Python.h>
*/
import "C"
import "unsafe"

type Str struct {
	Object
}

func newStr(obj *cPyObject) Str {
	return Str{newObject(obj)}
}

func MakeStr(s string) Str {
	ptr := (*Char)(unsafe.Pointer(unsafe.StringData(s)))
	length := C.long(len(s))
	return newStr(C.PyUnicode_FromStringAndSize(ptr, length))
}

func (s Str) String() string {
	var l C.long
	buf := C.PyUnicode_AsUTF8AndSize(s.obj, &l)
	return GoStringN((*Char)(buf), int(l))
}

func (s Str) Len() int {
	return int(C.PyUnicode_GetLength(s.obj))
}

func (s Str) ByteLen() int {
	var l C.long
	_ = C.PyUnicode_AsUTF8AndSize(s.obj, &l)
	return int(l)
}

func (s Str) Encode(encoding string) Bytes {
	return cast[Bytes](s.Call("encode", MakeStr(encoding)))
}
