package gp

/*
#include <Python.h>
*/
import "C"
import (
	"unsafe"
)

type Bytes struct {
	Object
}

func newBytes(obj *PyObject) Bytes {
	return Bytes{newObject(obj)}
}

func BytesFromStr(s string) Bytes {
	return MakeBytes([]byte(s))
}

func MakeBytes(bytes []byte) Bytes {
	ptr := C.CBytes(bytes)
	return newBytes(C.PyBytes_FromStringAndSize((*C.char)(ptr), C.Py_ssize_t(len(bytes))))
}

func (b Bytes) Bytes() []byte {
	var p *byte
	var l int
	return C.GoBytes(unsafe.Pointer(p), C.int(l))
}

func (b Bytes) Decode(encoding string) Str {
	return Cast[Str](b.Call("decode", MakeStr(encoding)))
}
