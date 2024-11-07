package gp

/*
#include <Python.h>
*/
import "C"
import "unsafe"

type Long struct {
	Object
}

func newLong(obj *cPyObject) Long {
	return Long{newObject(obj)}
}

func MakeLong(i int64) Long {
	return newLong(C.PyLong_FromLongLong(C.longlong(i)))
}

func LongFromFloat64(v float64) Long {
	return newLong(C.PyLong_FromDouble(C.double(v)))
}

func LongFromString(s string, base int) Long {
	cstr := AllocCStr(s)
	o := C.PyLong_FromString(cstr, nil, C.int(base))
	C.free(unsafe.Pointer(cstr))
	return newLong(o)
}

func LongFromUnicode(u Object, base int) Long {
	return newLong(C.PyLong_FromUnicodeObject(u.cpyObj(), C.int(base)))
}

func (l Long) Int() int {
	return int(l.Int64())
}

func (l Long) Int64() int64 {
	return int64(C.PyLong_AsLongLong(l.obj))
}

func (l Long) Uint() uint {
	return uint(l.Uint64())
}

func (l Long) Uint64() uint64 {
	return uint64(C.PyLong_AsUnsignedLongLong(l.obj))
}

func (l Long) Uintptr() uintptr {
	return uintptr(l.Int64())
}
func (l Long) Float64() float64 {
	return float64(C.PyLong_AsDouble(l.obj))
}

func LongFromUintptr(v uintptr) Long {
	return newLong(C.PyLong_FromLong(C.long(v)))
}
