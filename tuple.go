package gp

/*
#include <Python.h>
*/
import "C"

type Tuple struct {
	Object
}

func newTuple(obj *PyObject) Tuple {
	return Tuple{newObject(obj)}
}

func MakeTupleWithLen(len int) Tuple {
	return newTuple(C.PyTuple_New(C.Py_ssize_t(len)))
}

func MakeTuple(args ...any) Tuple {
	tuple := newTuple(C.PyTuple_New(C.Py_ssize_t(len(args))))
	for i, arg := range args {
		obj := From(arg)
		tuple.Set(i, obj)
	}
	return tuple
}

func (t Tuple) Get(index int) Object {
	v := C.PyTuple_GetItem(t.obj, C.Py_ssize_t(index))
	C.Py_IncRef(v)
	return newObject(v)
}

func (t Tuple) Set(index int, obj Objecter) {
	C.PyTuple_SetItem(t.obj, C.Py_ssize_t(index), obj.Obj())
}

func (t Tuple) Len() int {
	return int(C.PyTuple_Size(t.obj))
}

func (t Tuple) Slice(low, high int) Tuple {
	return newTuple(C.PyTuple_GetSlice(t.obj, C.Py_ssize_t(low), C.Py_ssize_t(high)))
}

func (t Tuple) ParseArgs(addrs ...any) bool {
	if len(addrs) > t.Len() {
		return false
	}

	for i, addr := range addrs {
		obj := t.Get(i)

		switch v := addr.(type) {
		// Integer types
		case *int:
			*v = int(obj.AsLong().Int64())
		case *int8:
			*v = int8(obj.AsLong().Int64())
		case *int16:
			*v = int16(obj.AsLong().Int64())
		case *int32:
			*v = int32(obj.AsLong().Int64())
		case *int64:
			*v = obj.AsLong().Int64()
		case *uint:
			*v = uint(obj.AsLong().Int64())
		case *uint8:
			*v = uint8(obj.AsLong().Int64())
		case *uint16:
			*v = uint16(obj.AsLong().Int64())
		case *uint32:
			*v = uint32(obj.AsLong().Int64())
		case *uint64:
			*v = uint64(obj.AsLong().Int64())

		// Floating point types
		case *float32:
			*v = float32(obj.AsFloat().Float64())
		case *float64:
			*v = obj.AsFloat().Float64()

		// Complex number types
		case *complex64:
			*v = complex64(obj.AsComplex().Complex128())
		case *complex128:
			*v = obj.AsComplex().Complex128()

		// String types
		case *string:
			*v = obj.AsStr().String()
		case *[]byte:
			*v = []byte(obj.AsStr().String())

		// Boolean type
		case *bool:
			*v = obj.AsBool().Bool()

		case **PyObject:
			*v = obj.Obj()

		// Python object
		case *Object:
			*v = obj

		default:
			return false
		}
	}

	return true
}
