package python

/*
#cgo pkg-config: python-3.12-embed
#include <Python.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type Str struct {
	Object
}

func newStr(obj *PyObject) Str {
	return Str{newObject(obj)}
}

func MakeStr(s string) Str {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	ptr := (*C.char)(unsafe.Pointer(hdr.Data))
	length := C.long(hdr.Len)
	return newStr(C.PyUnicode_FromStringAndSize(ptr, length))
}

func (s Str) String() string {
	var l C.long
	buf := C.PyUnicode_AsUTF8AndSize(s.obj, &l)
	return GoStringN((*C.char)(buf), int(l))
}

func (s Str) Len() int {
	var l C.long
	C.PyUnicode_AsUTF8AndSize(s.obj, &l)
	return int(l)
}

func (s Str) Encode(encoding string) Bytes {
	return Cast[Bytes](s.CallMethod("encode", MakeStr(encoding)))
}
