package gp

/*
#include <Python.h>
*/
import "C"

// Float represents a Python float object. It provides methods to convert between
// Go float types and Python float objects, as well as checking numeric properties.
type Float struct {
	Object
}

func newFloat(obj *cPyObject) Float {
	return Float{newObject(obj)}
}

func MakeFloat(f float64) Float {
	return newFloat(C.PyFloat_FromDouble(C.double(f)))
}

func (f Float) Float64() float64 {
	return float64(C.PyFloat_AsDouble(f.obj))
}

func (f Float) Float32() float32 {
	return float32(C.PyFloat_AsDouble(f.obj))
}

func (f Float) IsInteger() Bool {
	fn := cast[Func](f.Attr("is_integer"))
	return cast[Bool](fn.callNoArgs())
}
