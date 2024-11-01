package gp

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"
import "unsafe"

type PyObject = C.PyObject
type PyCFunction = C.PyCFunction

func Initialize() {
	C.Py_Initialize()
}

func Finalize() {
	mainMod = Module{}
	C.Py_FinalizeEx()
}

// ----------------------------------------------------------------------------

type InputType = C.int

const (
	SingleInput InputType = C.Py_single_input
	FileInput   InputType = C.Py_file_input
	EvalInput   InputType = C.Py_eval_input
)

func CompileString(code, filename string, start InputType) Object {
	return newObject(C.Py_CompileString(AllocCStr(code), AllocCStr(filename), C.int(start)))
}

func EvalCode(code Object, globals, locals Dict) Object {
	return newObject(C.PyEval_EvalCode(code.Obj(), globals.Obj(), locals.Obj()))
}

// ----------------------------------------------------------------------------

// llgo:link Cast llgo.staticCast
func Cast[U, T Objecter](obj T) (u U) {
	*(*T)(unsafe.Pointer(&u)) = obj
	return
}

// ----------------------------------------------------------------------------

func With[T Objecter](obj T, fn func(v T)) T {
	obj.object().Call("__enter__")
	defer obj.object().Call("__exit__")
	fn(obj)
	return obj
}

// ----------------------------------------------------------------------------

var mainMod Module

func MainModule() Module {
	if mainMod.Nil() {
		mainMod = ImportModule("__main__")
	}
	return mainMod
}

var noneObj Object

/*
from Dojo:
if self.none_value.is_null():

	var list_obj = self.PyList_New(0)
	var tuple_obj = self.PyTuple_New(0)
	var callable_obj = self.PyObject_GetAttrString(list_obj, "reverse")
	self.none_value = self.PyObject_CallObject(callable_obj, tuple_obj)
	self.Py_DecRef(tuple_obj)
	self.Py_DecRef(callable_obj)
	self.Py_DecRef(list_obj)
*/
func None() Object {
	if noneObj.Nil() {
		listObj := MakeList()
		tupleObj := MakeTuple()
		callableObj := listObj.GetFuncAttr("reverse")
		noneObj = callableObj.CallObject(tupleObj)
	}
	return noneObj
}

func Nil() Object {
	return Object{}
}
