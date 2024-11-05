package gp

/*
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// pyObject is a wrapper type that holds a Python Object and automatically calls
// the Python Object's DecRef method during garbage collection.
type pyObject struct {
	obj *C.PyObject
}

func (obj *pyObject) Obj() *PyObject {
	if obj == nil {
		return nil
	}
	return obj.obj
}

func (obj *pyObject) Nil() bool {
	return obj == nil
}

func (obj *pyObject) Ensure() {
	if obj == nil {
		C.PyErr_Print()
		panic("nil Python object")
	}
}

// ----------------------------------------------------------------------------

type Object struct {
	*pyObject
}

func FromPy(obj *PyObject) Object {
	return newObject(obj)
}

func (obj Object) object() Object {
	return obj
}

func newObject(obj *PyObject) Object {
	if obj == nil {
		C.PyErr_Print()
		panic("nil Python object")
	}
	o := &pyObject{obj: obj}
	p := Object{o}
	runtime.SetFinalizer(o, func(o *pyObject) {
		// TODO: need better auto-release mechanism
		// C.Py_DecRef(o.obj)
	})
	return p
}

func (obj Object) Dir() List {
	return obj.Call("__dir__").AsList()
}

func (obj Object) Equals(other Objecter) bool {
	return C.PyObject_RichCompareBool(obj.obj, other.Obj(), C.Py_EQ) != 0
}

func (obj Object) Attr(name string) Object {
	cname := AllocCStr(name)
	o := C.PyObject_GetAttrString(obj.obj, cname)
	C.free(unsafe.Pointer(cname))
	return newObject(o)
}

func (obj Object) AttrFloat(name string) Float {
	return obj.Attr(name).AsFloat()
}

func (obj Object) AttrLong(name string) Long {
	return obj.Attr(name).AsLong()
}

func (obj Object) AttrString(name string) Str {
	return obj.Attr(name).AsStr()
}

func (obj Object) AttrBytes(name string) Bytes {
	return obj.Attr(name).AsBytes()
}

func (obj Object) AttrBool(name string) Bool {
	return obj.Attr(name).AsBool()
}

func (obj Object) AttrDict(name string) Dict {
	return obj.Attr(name).AsDict()
}

func (obj Object) AttrList(name string) List {
	return obj.Attr(name).AsList()
}

func (obj Object) AttrTuple(name string) Tuple {
	return obj.Attr(name).AsTuple()
}

func (obj Object) AttrFunc(name string) Func {
	return obj.Attr(name).AsFunc()
}

func (obj Object) SetAttr(name string, value any) {
	cname := AllocCStr(name)
	r := C.PyObject_SetAttrString(obj.obj, cname, From(value).obj)
	C.PyErr_Print()
	check(r == 0, fmt.Sprintf("failed to set attribute %s", name))
	C.free(unsafe.Pointer(cname))
}

func (obj Object) IsLong() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyLong_Type) != 0
}

func (obj Object) IsFloat() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyFloat_Type) != 0
}

func (obj Object) IsComplex() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyComplex_Type) != 0
}

func (obj Object) IsStr() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyUnicode_Type) != 0
}

func (obj Object) IsBytes() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyBytes_Type) != 0
}

func (obj Object) IsBool() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyBool_Type) != 0
}

func (obj Object) IsList() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyList_Type) != 0
}

func (obj Object) IsTuple() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyTuple_Type) != 0
}

func (obj Object) IsDict() bool {
	return C.Py_IS_TYPE(obj.obj, &C.PyDict_Type) != 0
}

func (obj Object) AsFloat() Float {
	return Cast[Float](obj)
}

func (obj Object) AsLong() Long {
	return Cast[Long](obj)
}

func (obj Object) AsComplex() Complex {
	return Cast[Complex](obj)
}

func (obj Object) AsStr() Str {
	return Cast[Str](obj)
}

func (obj Object) AsBytes() Bytes {
	return Cast[Bytes](obj)
}

func (obj Object) AsBool() Bool {
	return Cast[Bool](obj)
}

func (obj Object) AsDict() Dict {
	return Cast[Dict](obj)
}

func (obj Object) AsList() List {
	return Cast[List](obj)
}

func (obj Object) AsTuple() Tuple {
	return Cast[Tuple](obj)
}

func (obj Object) AsFunc() Func {
	return Cast[Func](obj)
}

func (obj Object) AsModule() Module {
	return Cast[Module](obj)
}

func (obj Object) Call(name string, args ...any) Object {
	fn := Cast[Func](obj.Attr(name))
	argsTuple, kwArgs := splitArgs(args...)
	if kwArgs == nil {
		return fn.CallObject(argsTuple)
	} else {
		return fn.CallObjectKw(argsTuple, kwArgs)
	}
}

func (obj Object) Repr() string {
	return newStr(C.PyObject_Repr(obj.obj)).String()
}

func (obj Object) Type() Object {
	return newObject(C.PyObject_Type(obj.Obj()))
}

func (obj Object) String() string {
	return newStr(C.PyObject_Str(obj.obj)).String()
}

func (obj Object) Obj() *PyObject {
	if obj.Nil() {
		return nil
	}
	return obj.pyObject.obj
}
