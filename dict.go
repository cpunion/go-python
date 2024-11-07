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

func newDict(obj *cPyObject) Dict {
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
	v := C.PyDict_GetItem(d.obj, key.cpyObj())
	C.Py_IncRef(v)
	return newObject(v)
}

func (d Dict) Set(key, value Objecter) {
	keyObj := key.cpyObj()
	valueObj := value.cpyObj()
	C.PyDict_SetItem(d.obj, keyObj, valueObj)
}

func (d Dict) SetString(key string, value Objecter) {
	valueObj := value.cpyObj()
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
	C.PyDict_DelItem(d.obj, key.cpyObj())
}

func (d Dict) Iter() *DictIter {
	return &DictIter{dict: d, pos: 0}
}

type DictIter struct {
	dict Dict
	pos  C.long
}

func (d *DictIter) HasNext() bool {
	pos := d.pos
	return C.PyDict_Next(d.dict.obj, &pos, nil, nil) != 0
}

func (d *DictIter) Next() (Object, Object) {
	var key, value *C.PyObject
	if C.PyDict_Next(d.dict.obj, &d.pos, &key, &value) == 0 {
		return Nil(), Nil()
	}
	C.Py_IncRef(key)
	C.Py_IncRef(value)
	return newObject(key), newObject(value)
}
