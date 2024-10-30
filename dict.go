package gp

/*
#include <Python.h>
*/
import "C"

type Dict struct {
	Object
}

func newDict(obj *PyObject) Dict {
	return Dict{newObject(obj)}
}

func NewDict(obj *PyObject) Dict {
	return newDict(obj)
}

func DictFromPairs(pairs ...any) Dict {
	if len(pairs)%2 != 0 {
		panic("DictFromPairs requires an even number of arguments")
	}
	dict := newDict(C.PyDict_New())
	for i := 0; i < len(pairs); i += 2 {
		key := From(pairs[i])
		value := From(pairs[i+1])
		dict.Set(key, value)
	}
	return dict
}

func MakeDict(m map[any]any) Dict {
	dict := newDict(C.PyDict_New())
	for key, value := range m {
		keyObj := From(key)
		valueObj := From(value)
		dict.Set(keyObj, valueObj)
	}
	return dict
}

func (d Dict) Get(key Objecter) Object {
	v := C.PyDict_GetItem(d.obj, key.Obj())
	C.Py_IncRef(v)
	return newObject(v)
}

func (d Dict) Set(key, value Object) {
	C.Py_IncRef(key.obj)
	C.Py_IncRef(value.obj)
	C.PyDict_SetItem(d.obj, key.obj, value.obj)
}

func (d Dict) SetString(key string, value Object) {
	C.Py_IncRef(value.obj)
	C.PyDict_SetItemString(d.obj, AllocCStr(key), value.obj)
}

func (d Dict) GetString(key string) Object {
	v := C.PyDict_GetItemString(d.obj, AllocCStr(key))
	C.Py_IncRef(v)
	return newObject(v)
}

func (d Dict) Del(key Object) {
	C.PyDict_DelItem(d.obj, key.obj)
	C.Py_DecRef(key.obj)
}
