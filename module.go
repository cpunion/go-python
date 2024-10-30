package gp

/*
#include <Python.h>
*/
import "C"

import (
	"unsafe"
)

type Module struct {
	Object
}

func newModule(obj *PyObject) Module {
	return Module{newObject(obj)}
}

func ImportModule(name string) Module {
	mod := C.PyImport_ImportModule(AllocCStr(name))
	return newModule(mod)
}

func (m Module) Dict() Dict {
	return newDict(C.PyModule_GetDict(m.obj))
}

func (m Module) AddObject(name string, obj Object) int {
	return int(C.PyModule_AddObject(m.obj, AllocCStr(name), obj.obj))
}

func (m Module) AddFunction(name string, fn unsafe.Pointer, doc string) Func {
	def := &C.PyMethodDef{
		ml_name:  AllocCStr(name),
		ml_meth:  C.PyCFunction(fn),
		ml_flags: C.METH_VARARGS,
		ml_doc:   AllocCStr(doc),
	}
	pyFn := C.PyCMethod_New(def, nil, m.obj, nil)
	return newFunc(pyFn)
}
