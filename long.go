package python

/*
#cgo pkg-config: python-3.12-embed
#include <Python.h>
*/
import "C"

type Long struct {
	Object
}

func newLong(obj *PyObject) Long {
	return Long{newObject(obj)}
}

func MakeLong(i int64) Long {
	return newLong(C.PyLong_FromLongLong(C.longlong(i)))
}

func (l Long) Int64() int64 {
	return int64(C.PyLong_AsLongLong(l.obj))
}

func (l Long) Uint64() uint64 {
	return uint64(C.PyLong_AsUnsignedLongLong(l.obj))
}

func (l Long) AsFloat64() float64 {
	return float64(C.PyLong_AsDouble(l.obj))
}

func LongFromFloat64(v float64) Long {
	return newLong(C.PyLong_FromDouble(C.double(v)))
}

func LongFromString(s string, base int) Long {
	cstr := AllocCStr(s)
	return newLong(C.PyLong_FromString(cstr, nil, C.int(base)))
}

func LongFromUnicode(u Object, base int) Long {
	return newLong(C.PyLong_FromUnicodeObject(u.Obj(), C.int(base)))
}

func (l Long) AsUint64() uint64 {
	return uint64(C.PyLong_AsUnsignedLongLong(l.obj))
}

func (l Long) AsUintptr() uintptr {
	return uintptr(C.PyLong_AsLong(l.obj))
}

func LongFromUintptr(v uintptr) Long {
	return newLong(C.PyLong_FromLong(C.long(v)))
}
