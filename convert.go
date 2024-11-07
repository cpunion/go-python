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

func From(from any) Object {
	switch v := from.(type) {
	case Objecter:
		return newObject(v.cpyObj())
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

func ToValue(from Object, to reflect.Value) bool {
	if !to.IsValid() || !to.CanSet() {
		panic(fmt.Errorf("value is not valid or cannot be set: %v\n", to))
	}

	switch to.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if from.IsLong() {
			to.SetInt(cast[Long](from).Int64())
		} else {
			return false
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if from.IsLong() {
			to.SetUint(cast[Long](from).Uint64())
		} else {
			return false
		}
	case reflect.Float32, reflect.Float64:
		if from.IsFloat() || from.IsLong() {
			to.SetFloat(cast[Float](from).Float64())
		} else {
			return false
		}
	case reflect.Complex64, reflect.Complex128:
		if from.IsComplex() {
			to.SetComplex(cast[Complex](from).Complex128())
		} else {
			return false
		}
	case reflect.String:
		if from.IsStr() {
			to.SetString(cast[Str](from).String())
		} else {
			return false
		}
	case reflect.Bool:
		if from.IsBool() {
			to.SetBool(cast[Bool](from).Bool())
		} else {
			return false
		}
	case reflect.Slice:
		if to.Type().Elem().Kind() == reflect.Uint8 { // []byte
			if from.IsBytes() {
				to.SetBytes(cast[Bytes](from).Bytes())
			} else {
				return false
			}
		} else {
			if from.IsList() {
				list := cast[List](from)
				l := list.Len()
				slice := reflect.MakeSlice(to.Type(), l, l)
				for i := 0; i < l; i++ {
					item := list.GetItem(i)
					ToValue(item, slice.Index(i))
				}
				to.Set(slice)
			} else {
				return false
			}
		}
	case reflect.Map:
		if from.IsDict() {
			t := to.Type()
			to.Set(reflect.MakeMap(t))
			dict := cast[Dict](from)
			iter := dict.Iter()
			for iter.HasNext() {
				key, value := iter.Next()
				vk := reflect.New(t.Key()).Elem()
				vv := reflect.New(t.Elem()).Elem()
				if !ToValue(key, vk) || !ToValue(value, vv) {
					return false
				}
				to.SetMapIndex(vk, vv)
			}
			return true
		} else {
			return false
		}
	case reflect.Struct:
		if from.IsDict() {
			dict := cast[Dict](from)
			t := to.Type()
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				key := goNameToPythonName(field.Name)
				if !dict.HasKey(MakeStr(key)) {
					continue
				}
				value := dict.Get(MakeStr(key))
				if !ToValue(value, to.Field(i)) {
					SetTypeError(fmt.Errorf("failed to convert value to %v", field.Name))
					return false
				}
			}
		} else {
			maps := getGlobalData()
			tyMeta := maps.typeMetas[from.Type().cpyObj()]
			if tyMeta == nil {
				return false
			}
			wrapper := (*wrapperType)(unsafe.Pointer(from.cpyObj()))
			to.Set(reflect.ValueOf(wrapper.goObj).Elem())
			return true
		}
	default:
		panic(fmt.Errorf("unsupported type conversion from Python object to %v", to.Type()))
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
