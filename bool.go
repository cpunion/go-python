package gp

/*
#include <Python.h>
*/
import "C"

type Bool struct {
	Object
}

func newBool(obj *PyObject) Bool {
	return Bool{newObject(obj)}
}

func MakeBool(b bool) Bool {
	if b {
		return True()
	}
	return False()
}

func True() Bool {
	return newBool(C.Py_True).AsBool()
}

func False() Bool {
	return newBool(C.Py_False).AsBool()
}

func (b Bool) Bool() bool {
	return C.PyObject_IsTrue(b.obj) != 0
}

func (b Bool) Not() Bool {
	return MakeBool(!b.Bool())
}
