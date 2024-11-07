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

func newBytes(obj *cPyObject) Bytes {
	return Bytes{newObject(obj)}
}

func BytesFromStr(s string) Bytes {
	return MakeBytes([]byte(s))
}

func MakeBytes(bytes []byte) Bytes {
	ptr := C.CBytes(bytes)
	o := C.PyBytes_FromStringAndSize((*Char)(ptr), C.Py_ssize_t(len(bytes)))
	C.free(unsafe.Pointer(ptr))
	return newBytes(o)
}

func (b Bytes) Bytes() []byte {
	p := (*byte)(unsafe.Pointer(C.PyBytes_AsString(b.obj)))
	l := int(C.PyBytes_Size(b.obj))
	return C.GoBytes(unsafe.Pointer(p), C.int(l))
}

func (b Bytes) Decode(encoding string) Str {
	return cast[Str](b.Call("decode", MakeStr(encoding)))
}
