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

var trueObj Object
var falseObj Object

func True() Bool {
	if trueObj.Nil() {
		trueObj = MainModule().Dict().GetString("True")
	}
	return trueObj.AsBool()
}

func False() Bool {
	if falseObj.Nil() {
		falseObj = MainModule().Dict().GetString("False")
	}
	return falseObj.AsBool()
}

func (b Bool) Bool() bool {
	return C.PyObject_IsTrue(b.obj) != 0
}
