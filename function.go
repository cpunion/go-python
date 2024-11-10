package gp

/*
#include <Python.h>
*/
import "C"

type Objecter interface {
	cpyObj() *cPyObject
	object() Object
	Ensure()
}

type Func struct {
	Object
}

func newFunc(obj *cPyObject) Func {
	return Func{newObject(obj)}
}

func (f Func) Ensure() {
	f.pyObject.Ensure()
}

func (f Func) Name() string {
	return f.AttrString("__name__").String()
}

func (f Func) call(args Tuple, kwargs Dict) Object {
	return newObject(C.PyObject_Call(f.obj, args.obj, kwargs.obj))
}

func (f Func) callNoArgs() Object {
	return newObject(C.PyObject_CallNoArgs(f.obj))
}

func (f Func) callOneArg(arg Objecter) Object {
	return newObject(C.PyObject_CallOneArg(f.obj, arg.cpyObj()))
}

func (f Func) CallObject(args Tuple) Object {
	defer getGlobalData().decRefObjectsIfNeeded()
	return newObject(C.PyObject_CallObject(f.obj, args.obj))
}

func (f Func) CallObjectKw(args Tuple, kw KwArgs) Object {
	defer getGlobalData().decRefObjectsIfNeeded()
	// Convert keyword arguments to Python dict
	kwDict := MakeDict(nil)
	for k, v := range kw {
		kwDict.Set(MakeStr(k), From(v))
	}
	return f.call(args, kwDict)
}

func (f Func) Call(args ...any) Object {
	argsTuple, kwArgs := splitArgs(args...)
	if kwArgs == nil {
		defer getGlobalData().decRefObjectsIfNeeded()
		switch len(args) {
		case 0:
			return f.callNoArgs()
		case 1:
			return f.callOneArg(From(args[0]))
		default:
			return f.CallObject(argsTuple)
		}
	} else {
		return f.CallObjectKw(argsTuple, kwArgs)
	}
}

// ----------------------------------------------------------------------------
