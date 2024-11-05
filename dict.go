package gp

/*
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Dict struct {
	Object
}

func newDict(obj *PyObject) Dict {
	return Dict{newObject(obj)}
}

func DictFromPairs(pairs ...any) Dict {
	check(len(pairs)%2 == 0, "DictFromPairs requires an even number of arguments")
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

func (d Dict) HasKey(key any) bool {
	keyObj := From(key)
	return C.PyDict_Contains(d.obj, keyObj.obj) != 0
}

func (d Dict) Get(key Objecter) Object {
	v := C.PyDict_GetItem(d.obj, key.Obj())
	C.Py_IncRef(v)
	return newObject(v)
}

func (d Dict) Set(key, value Objecter) {
	keyObj := key.Obj()
	valueObj := value.Obj()
	C.PyDict_SetItem(d.obj, keyObj, valueObj)
}

func (d Dict) SetString(key string, value Objecter) {
	valueObj := value.Obj()
	ckey := AllocCStr(key)
	r := C.PyDict_SetItemString(d.obj, ckey, valueObj)
	C.free(unsafe.Pointer(ckey))
	check(r == 0, fmt.Sprintf("failed to set item string: %v", r))
}

func (d Dict) GetString(key string) Object {
	ckey := AllocCStr(key)
	v := C.PyDict_GetItemString(d.obj, ckey)
	C.Py_IncRef(v)
	C.free(unsafe.Pointer(ckey))
	return newObject(v)
}

func (d Dict) Del(key Objecter) {
	C.PyDict_DelItem(d.obj, key.Obj())
}

func (d Dict) Items() func(fn func(key, value Object) bool) {
	return func(fn func(key, value Object) bool) {
		items := C.PyDict_Items(d.obj)
		check(items != nil, "failed to get items of dict")
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
			C.Py_IncRef(key)
			C.Py_IncRef(value)
			C.Py_DecRef(item)
			if !fn(newObject(key), newObject(value)) {
				break
			}
		}
	}
}
