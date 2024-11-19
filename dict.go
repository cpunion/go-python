package gp

/*
#include <Python.h>

typedef struct pyCriticalSection {
    uintptr_t _cs_prev;
    void *_cs_mutex;
} pyCriticalSection;
static inline void pyCriticalSection_Begin(pyCriticalSection *pcs, PyObject *op) {
#if PY_VERSION_HEX >= 0x030D0000
    PyCriticalSection_Begin((PyCriticalSection*)pcs, op);
#else
    PyGILState_STATE gstate = PyGILState_Ensure();
		pcs->_cs_prev = (uintptr_t)gstate;
#endif
}
static inline void pyCriticalSection_End(pyCriticalSection *pcs) {
#if PY_VERSION_HEX >= 0x030D0000
    PyCriticalSection_End((PyCriticalSection*)pcs);
#else
    PyGILState_Release((PyGILState_STATE)pcs->_cs_prev);
#endif
}
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

func (d Dict) Items() func(func(Object, Object) bool) {
	obj := d.cpyObj()
	var cs C.pyCriticalSection
	C.pyCriticalSection_Begin(&cs, obj)
	return func(fn func(Object, Object) bool) {
		defer C.pyCriticalSection_End(&cs)
		var pos C.Py_ssize_t
		var key, value *C.PyObject
		for C.PyDict_Next(obj, &pos, &key, &value) == 1 {
			if !fn(newObject(key), newObject(value)) {
				return
			}
		}
	}
}
