package gp

/*
#include <Python.h>
*/
import "C"
import "fmt"

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

func (d Dict) ForEach(fn func(key, value Object)) {
	items := C.PyDict_Items(d.obj)
	if items == nil {
		panic(fmt.Errorf("failed to get items of dict"))
	}
	defer C.Py_DecRef(items)
	iter := C.PyObject_GetIter(items)
	for {
		item := C.PyIter_Next(iter)
		if item == nil {
			break
		}
		C.Py_IncRef(item)
		key := C.PyTuple_GetItem(item, 0)
		value := C.PyTuple_GetItem(item, 1)
		fn(newObject(key), newObject(value))
	}
}
