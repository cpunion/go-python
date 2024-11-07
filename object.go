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
	g   *globalData
}

func (o *pyObject) cpyObj() *cPyObject {
	if o == nil {
		return nil
	}
	return o.obj
}

func (o *pyObject) Nil() bool {
	return o == nil
}

func (o *pyObject) Ensure() {
	if o == nil {
		C.PyErr_Print()
		panic("nil Python object")
	}
}

// ----------------------------------------------------------------------------

type Object struct {
	*pyObject
}

func FromPy(obj *cPyObject) Object {
	return newObject(obj)
}

func (o Object) object() Object {
	return o
}

func newObjectRef(obj *cPyObject) Object {
	C.Py_IncRef(obj)
	return newObject(obj)
}

func newObject(obj *cPyObject) Object {
	if obj == nil {
		C.PyErr_Print()
		panic("nil Python object")
	}
	o := &pyObject{obj: obj, g: getGlobalData()}
	runtime.SetFinalizer(o, func(o *pyObject) {
		o.g.addDecRef(o.obj)
		runtime.SetFinalizer(o, nil)
	})
	return Object{o}
}

func (o Object) Dir() List {
	return o.Call("__dir__").AsList()
}

func (o Object) Equals(other Objecter) bool {
	return C.PyObject_RichCompareBool(o.obj, other.cpyObj(), C.Py_EQ) != 0
}

func (o Object) Attr(name string) Object {
	cname := AllocCStr(name)
	attr := C.PyObject_GetAttrString(o.obj, cname)
	C.free(unsafe.Pointer(cname))
	return newObject(attr)
}

func (o Object) AttrFloat(name string) Float {
	return o.Attr(name).AsFloat()
}

func (o Object) AttrLong(name string) Long {
	return o.Attr(name).AsLong()
}

func (o Object) AttrString(name string) Str {
	return o.Attr(name).AsStr()
}

func (o Object) AttrBytes(name string) Bytes {
	return o.Attr(name).AsBytes()
}

func (o Object) AttrBool(name string) Bool {
	return o.Attr(name).AsBool()
}

func (o Object) AttrDict(name string) Dict {
	return o.Attr(name).AsDict()
}

func (o Object) AttrList(name string) List {
	return o.Attr(name).AsList()
}

func (o Object) AttrTuple(name string) Tuple {
	return o.Attr(name).AsTuple()
}

func (o Object) AttrFunc(name string) Func {
	return o.Attr(name).AsFunc()
}

func (o Object) SetAttr(name string, value any) {
	cname := AllocCStr(name)
	r := C.PyObject_SetAttrString(o.obj, cname, From(value).obj)
	C.PyErr_Print()
	check(r == 0, fmt.Sprintf("failed to set attribute %s", name))
	C.free(unsafe.Pointer(cname))
}

func (o Object) IsLong() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyLong_Type) != 0
}

func (o Object) IsFloat() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyFloat_Type) != 0
}

func (o Object) IsComplex() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyComplex_Type) != 0
}

func (o Object) IsStr() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyUnicode_Type) != 0
}

func (o Object) IsBytes() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyBytes_Type) != 0
}

func (o Object) IsBool() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyBool_Type) != 0
}

func (o Object) IsList() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyList_Type) != 0
}

func (o Object) IsTuple() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyTuple_Type) != 0
}

func (o Object) IsDict() bool {
	return C.Py_IS_TYPE(o.obj, &C.PyDict_Type) != 0
}

func (o Object) AsFloat() Float {
	return cast[Float](o)
}

func (o Object) AsLong() Long {
	return cast[Long](o)
}

func (o Object) AsComplex() Complex {
	return cast[Complex](o)
}

func (o Object) AsStr() Str {
	return cast[Str](o)
}

func (o Object) AsBytes() Bytes {
	return cast[Bytes](o)
}

func (o Object) AsBool() Bool {
	return cast[Bool](o)
}

func (o Object) AsDict() Dict {
	return cast[Dict](o)
}

func (o Object) AsList() List {
	return cast[List](o)
}

func (o Object) AsTuple() Tuple {
	return cast[Tuple](o)
}

func (o Object) AsFunc() Func {
	return cast[Func](o)
}

func (o Object) AsModule() Module {
	return cast[Module](o)
}

func (o Object) Call(name string, args ...any) Object {
	fn := cast[Func](o.Attr(name))
	argsTuple, kwArgs := splitArgs(args...)
	if kwArgs == nil {
		return fn.CallObject(argsTuple)
	} else {
		return fn.CallObjectKw(argsTuple, kwArgs)
	}
}

func (o Object) Repr() string {
	return newStr(C.PyObject_Repr(o.obj)).String()
}

func (o Object) Type() Object {
	return newObject(C.PyObject_Type(o.cpyObj()))
}

func (o Object) String() string {
	return newStr(C.PyObject_Str(o.obj)).String()
}

func (o Object) cpyObj() *cPyObject {
	if o.Nil() {
		return nil
	}
	return o.obj
}
