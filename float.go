package python

/*
#cgo pkg-config: python-3.12-embed
#include <Python.h>
*/
import "C"

type Float struct {
	Object
}

func newFloat(obj *C.PyObject) Float {
	return Float{newObject(obj)}
}

func MakeFloat(f float64) Float {
	return newFloat(C.PyFloat_FromDouble(C.double(f)))
}

func (f Float) Float64() float64 {
	return float64(C.PyFloat_AsDouble(f.obj))
}

func (f Float) IsInteger() Bool {
	fn := Cast[Func](f.GetAttr("is_integer"))
	return Cast[Bool](fn.callNoArgs())
}
