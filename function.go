package gp

/*
#include <Python.h>
*/
import "C"

type Objecter interface {
	Obj() *PyObject
	object() Object
	Ensure()
}

type Func struct {
	Object
}

func newFunc(obj *PyObject) Func {
	return Func{newObject(obj)}
}

func (f Func) Ensure() {
	f.pyObject.Ensure()
}

func (f Func) call(args Tuple, kwargs Dict) Object {
	return newObject(C.PyObject_Call(f.obj, args.obj, kwargs.obj))
}

func (f Func) callNoArgs() Object {
	return newObject(C.PyObject_CallNoArgs(f.obj))
}

func (f Func) callOneArg(arg Objecter) Object {
	return newObject(C.PyObject_CallOneArg(f.obj, arg.Obj()))
}

func (f Func) CallObject(args Tuple) Object {
	return newObject(C.PyObject_CallObject(f.obj, args.obj))
}

func (f Func) CallObjectKw(args Tuple, kw KwArgs) Object {
	return f.call(args, From(map[string]any(kw)).AsDict())
}

func (f Func) Call(args ...any) Object {
	switch len(args) {
	case 0:
		return f.callNoArgs()
	case 1:
		return f.callOneArg(From(args[0]))
	default:
		argsTuple := C.PyTuple_New(C.Py_ssize_t(len(args)))
		for i, arg := range args {
			obj := From(arg).Obj()
			C.Py_IncRef(obj)
			C.PyTuple_SetItem(argsTuple, C.Py_ssize_t(i), obj)
		}
		return newObject(C.PyObject_CallObject(f.obj, argsTuple))
	}
}

// ----------------------------------------------------------------------------
