package gp

/*
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
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
		return None()
	}
	o := &pyObject{obj: obj}
	p := Object{o}
	runtime.SetFinalizer(o, func(o *pyObject) {
		C.Py_DecRef(o.obj)
	})
	return p
}

func (obj Object) Dir() List {
	return obj.Call("__dir__").AsList()
}

func (obj Object) GetAttr(name string) Object {
	o := C.PyObject_GetAttrString(obj.obj, AllocCStr(name))
	C.Py_IncRef(o)
	return newObject(o)
}

func (obj Object) GetFloatAttr(name string) Float {
	return obj.GetAttr(name).AsFloat()
}

func (obj Object) GetLongAttr(name string) Long {
	return obj.GetAttr(name).AsLong()
}

func (obj Object) GetStrAttr(name string) Str {
	return obj.GetAttr(name).AsStr()
}

func (obj Object) GetBytesAttr(name string) Bytes {
	return obj.GetAttr(name).AsBytes()
}

func (obj Object) GetBoolAttr(name string) Bool {
	return obj.GetAttr(name).AsBool()
}

func (obj Object) GetDictAttr(name string) Dict {
	return obj.GetAttr(name).AsDict()
}

func (obj Object) GetListAttr(name string) List {
	return obj.GetAttr(name).AsList()
}

func (obj Object) GetTupleAttr(name string) Tuple {
	return obj.GetAttr(name).AsTuple()
}

func (obj Object) GetFuncAttr(name string) Func {
	return obj.GetAttr(name).AsFunc()
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
	fn := Cast[Func](obj.GetAttr(name))
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

func (obj Object) String() string {
	return newStr(C.PyObject_Str(obj.obj)).String()
}

func (obj Object) Obj() *PyObject {
	if obj.Nil() {
		return nil
	}
	return obj.pyObject.obj
}

func From(v any) Object {
	switch v := v.(type) {
	case Objecter:
		return newObject(v.Obj())
	case int8:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case int16:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case int32:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case int64:
		return newObject(C.PyLong_FromLongLong(C.longlong(v)))
	case int:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case uint8:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case uint16:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case uint32:
		return newObject(C.PyLong_FromLong(C.long(v)))
	case uint64:
		return newObject(C.PyLong_FromUnsignedLongLong(C.ulonglong(v)))
	case uint:
		return newObject(C.PyLong_FromUnsignedLong(C.ulong(v)))
	case float64:
		return newObject(C.PyFloat_FromDouble(C.double(v)))
	case string:
		return newObject(C.PyUnicode_FromString(AllocCStr(v)))
	case complex128:
		return MakeComplex(v).Object
	case complex64:
		return MakeComplex(complex128(v)).Object
	case []byte:
		return MakeBytes(v).Object
	case bool:
		if v {
			return True().Object
		} else {
			return False().Object
		}
	case *C.PyObject:
		return newObject(v)
	default:
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Slice:
			return fromSlice(vv).Object
		case reflect.Map:
			return fromMap(vv).Object
		}
		panic(fmt.Errorf("unsupported type for Python: %T\n", v))
	}
}

func ToValue(obj Object, v reflect.Value) bool {
	// Handle nil pointer
	if !v.IsValid() || !v.CanSet() {
		return false
	}

	switch v.Kind() {
	case reflect.Int8:
		v.SetInt(Cast[Long](obj).Int64())
	case reflect.Int16:
		v.SetInt(Cast[Long](obj).Int64())
	case reflect.Int32:
		v.SetInt(Cast[Long](obj).Int64())
	case reflect.Int64:
		v.SetInt(Cast[Long](obj).Int64())
	case reflect.Int:
		v.SetInt(Cast[Long](obj).Int64())
	case reflect.Uint8:
		v.SetUint(Cast[Long](obj).Uint64())
	case reflect.Uint16:
		v.SetUint(Cast[Long](obj).Uint64())
	case reflect.Uint32:
		v.SetUint(Cast[Long](obj).Uint64())
	case reflect.Uint64:
		v.SetUint(Cast[Long](obj).Uint64())
	case reflect.Uint:
		v.SetUint(Cast[Long](obj).Uint64())
	case reflect.Float32:
		v.SetFloat(Cast[Float](obj).Float64())
	case reflect.Float64:
		v.SetFloat(Cast[Float](obj).Float64())
	case reflect.Complex64, reflect.Complex128:
		v.SetComplex(Cast[Complex](obj).Complex128())
	case reflect.String:
		v.SetString(Cast[Str](obj).String())
	case reflect.Bool:
		v.SetBool(Cast[Bool](obj).Bool())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 { // []byte
			v.SetBytes(Cast[Bytes](obj).Bytes())
		} else {
			list := Cast[List](obj)
			l := list.Len()
			slice := reflect.MakeSlice(v.Type(), l, l)
			for i := 0; i < l; i++ {
				item := list.GetItem(i)
				ToValue(item, slice.Index(i))
			}
			v.Set(slice)
		}
	default:
		panic(fmt.Errorf("unsupported type conversion from Python object to %v", v.Type()))
	}
	return true
}

func To[T any](obj Object) (ret T) {
	switch any(ret).(type) {
	case int8:
		return any(int8(Cast[Long](obj).Int64())).(T)
	case int16:
		return any(int16(Cast[Long](obj).Int64())).(T)
	case int32:
		return any(int32(Cast[Long](obj).Int64())).(T)
	case int64:
		return any(Cast[Long](obj).Int64()).(T)
	case int:
		return any(int(Cast[Long](obj).Int64())).(T)
	case uint8:
		return any(uint8(Cast[Long](obj).Uint64())).(T)
	case uint16:
		return any(uint16(Cast[Long](obj).Uint64())).(T)
	case uint32:
		return any(uint32(Cast[Long](obj).Uint64())).(T)
	case uint64:
		return any(Cast[Long](obj).Uint64()).(T)
	case uint:
		return any(uint(Cast[Long](obj).Uint64())).(T)
	case float32:
		return any(float32(Cast[Float](obj).Float64())).(T)
	case float64:
		return any(Cast[Float](obj).Float64()).(T)
	case complex64:
		return any(complex64(Cast[Complex](obj).Complex128())).(T)
	case complex128:
		return any(Cast[Complex](obj).Complex128()).(T)
	case string:
		return any(Cast[Str](obj).String()).(T)
	case bool:
		return any(Cast[Bool](obj).Bool()).(T)
	case []byte:
		return any(Cast[Bytes](obj).Bytes()).(T)
	default:
		v := reflect.ValueOf(ret)
		switch v.Kind() {
		case reflect.Slice:
			return toSlice[T](obj, v)
		}
		panic(fmt.Errorf("unsupported type conversion from Python object to %T", ret))
	}
}

func toSlice[T any](obj Object, v reflect.Value) T {
	list := Cast[List](obj)
	l := list.Len()
	v = reflect.MakeSlice(v.Type(), l, l)
	for i := 0; i < l; i++ {
		v.Index(i).Set(reflect.ValueOf(To[T](list.GetItem(i))))
	}
	return v.Interface().(T)
}

func fromSlice(v reflect.Value) List {
	l := v.Len()
	list := newList(C.PyList_New(C.Py_ssize_t(l)))
	for i := 0; i < l; i++ {
		list.SetItem(i, From(v.Index(i).Interface()))
	}
	return list
}

func fromMap(v reflect.Value) Dict {
	dict := newDict(C.PyDict_New())
	for _, key := range v.MapKeys() {
		dict.Set(From(key.Interface()), From(v.MapIndex(key).Interface()))
	}
	return dict
}
