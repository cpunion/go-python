package gp

/*
#include <Python.h>
*/
import "C"

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

func GetModule(name string) Module {
	return newModule(C.PyImport_GetModule(MakeStr(name).obj))
}

func (m Module) Dict() Dict {
	return newDict(C.PyModule_GetDict(m.obj))
}

func (m Module) AddObject(name string, obj Object) int {
	return int(C.PyModule_AddObject(m.obj, AllocCStr(name), obj.obj))
}

func CreateModule(name string) Module {
	return newModule(C.PyModule_New(AllocCStr(name)))
}
