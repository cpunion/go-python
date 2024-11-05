package gp

/*
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

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
	case float32:
		return newObject(C.PyFloat_FromDouble(C.double(v)))
	case string:
		cstr := AllocCStr(v)
		o := C.PyUnicode_FromString(cstr)
		C.free(unsafe.Pointer(cstr))
		return newObject(o)
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
		case reflect.Ptr:
			if vv.Elem().Type().Kind() == reflect.Struct {
				maps := getGlobalData()
				if pyType, ok := maps.pyTypes[vv.Elem().Type()]; ok {
					wrapper := allocWrapper((*C.PyTypeObject)(unsafe.Pointer(pyType)), vv.Interface())
					return newObject((*C.PyObject)(unsafe.Pointer(wrapper)))
				}
			}
			return From(vv.Elem().Interface())
		case reflect.Slice:
			return fromSlice(vv).Object
		case reflect.Map:
			return fromMap(vv).Object
		case reflect.Struct:
			return fromStruct(vv)
		}
		panic(fmt.Errorf("unsupported type for Python: %T\n", v))
	}
}

func ToValue(obj Object, v reflect.Value) bool {
	if !v.IsValid() || !v.CanSet() {
		panic(fmt.Errorf("value is not valid or cannot be set: %v\n", v))
	}

	switch v.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if obj.IsLong() {
			v.SetInt(Cast[Long](obj).Int64())
		} else {
			return false
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if obj.IsLong() {
			v.SetUint(Cast[Long](obj).Uint64())
		} else {
			return false
		}
	case reflect.Float32, reflect.Float64:
		if obj.IsFloat() || obj.IsLong() {
			v.SetFloat(Cast[Float](obj).Float64())
		} else {
			return false
		}
	case reflect.Complex64, reflect.Complex128:
		if obj.IsComplex() {
			v.SetComplex(Cast[Complex](obj).Complex128())
		} else {
			return false
		}
	case reflect.String:
		if obj.IsStr() {
			v.SetString(Cast[Str](obj).String())
		} else {
			return false
		}
	case reflect.Bool:
		if obj.IsBool() {
			v.SetBool(Cast[Bool](obj).Bool())
		} else {
			return false
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 { // []byte
			if obj.IsBytes() {
				v.SetBytes(Cast[Bytes](obj).Bytes())
			} else {
				return false
			}
		} else {
			if obj.IsList() {
				list := Cast[List](obj)
				l := list.Len()
				slice := reflect.MakeSlice(v.Type(), l, l)
				for i := 0; i < l; i++ {
					item := list.GetItem(i)
					ToValue(item, slice.Index(i))
				}
				v.Set(slice)
			} else {
				return false
			}
		}
	case reflect.Map:
		if obj.IsDict() {
			t := v.Type()
			v.Set(reflect.MakeMap(t))
			dict := Cast[Dict](obj)
			dict.ForEach(func(key, value Object) {
				vk := reflect.New(t.Key()).Elem()
				vv := reflect.New(t.Elem()).Elem()
				if !ToValue(key, vk) || !ToValue(value, vv) {
					panic(fmt.Errorf("failed to convert key or value to %v", t.Key()))
				}
				v.SetMapIndex(vk, vv)
			})
		} else {
			return false
		}
	case reflect.Struct:
		if obj.IsDict() {
			dict := Cast[Dict](obj)
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				key := goNameToPythonName(field.Name)
				value := dict.Get(MakeStr(key))
				if !ToValue(value, v.Field(i)) {
					panic(fmt.Errorf("failed to convert value to %v", field.Name))
				}
			}
		} else {
			maps := getGlobalData()
			tyMeta := maps.typeMetas[obj.Type().Obj()]
			if tyMeta == nil {
				return false
			}
			wrapper := (*wrapperType)(unsafe.Pointer(obj.Obj()))
			v.Set(reflect.ValueOf(wrapper.goObj).Elem())
			return true
		}
	default:
		panic(fmt.Errorf("unsupported type conversion from Python object to %v", v.Type()))
	}
	return true
}

func fromSlice(v reflect.Value) List {
	l := v.Len()
	list := newList(C.PyList_New(C.Py_ssize_t(l)))
	ty := v.Type().Elem()
	maps := getGlobalData()
	pyType, ok := maps.pyTypes[ty]
	if !ok {
		for i := 0; i < l; i++ {
			C.PyList_SetItem(list.obj, C.Py_ssize_t(i), From(v.Index(i).Interface()).obj)
		}
	} else {
		for i := 0; i < l; i++ {
			elem := v.Index(i)
			elemAddr := elem.Addr()
			wrapper := allocWrapper((*C.PyTypeObject)(unsafe.Pointer(pyType)), 0)
			wrapper.goObj = elemAddr.Interface()
			C.PyList_SetItem(list.obj, C.Py_ssize_t(i), (*C.PyObject)(unsafe.Pointer(wrapper)))
		}
	}
	return list
}

func fromMap(v reflect.Value) Dict {
	dict := newDict(C.PyDict_New())
	iter := v.MapRange()
	for iter.Next() {
		dict.Set(From(iter.Key().Interface()), From(iter.Value().Interface()))
	}
	return dict
}

func fromStruct(v reflect.Value) Object {
	ty := v.Type()
	maps := getGlobalData()
	if typeObj, ok := maps.pyTypes[ty]; ok {
		ptr := reflect.New(ty)
		ptr.Elem().Set(v)
		wrapper := allocWrapper((*C.PyTypeObject)(unsafe.Pointer(typeObj)), ptr.Interface())
		return newObject((*C.PyObject)(unsafe.Pointer(wrapper)))
	}
	dict := newDict(C.PyDict_New())
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		key := goNameToPythonName(field.Name)
		dict.Set(MakeStr(key).Object, From(v.Field(i).Interface()))
	}
	return dict.Object
}
