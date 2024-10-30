package python

/*
#cgo pkg-config: python-3.12-embed
#include <Python.h>

extern PyObject* wrapperFunc(PyObject* self, PyObject* args);
*/
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

type Objecter interface {
	Obj() *C.PyObject
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

type wrapperContext struct {
	v any
	t reflect.Type
}

//export wrapperFunc
func wrapperFunc(self, args *PyObject) *PyObject {
	wCtx := (*wrapperContext)(C.PyCapsule_GetPointer(self, AllocCStr("wrapperContext")))
	v := reflect.ValueOf(wCtx.v)
	t := v.Type()

	goArgs := make([]reflect.Value, t.NumIn())
	for i := range goArgs {
		goArgs[i] = reflect.New(t.In(i)).Elem()
		ToValue(FromPy(C.PyTuple_GetItem(args, C.Py_ssize_t(i))), goArgs[i])
	}

	results := v.Call(goArgs)

	if len(results) == 0 {
		return None().Obj()
	}
	if len(results) == 1 {
		return From(results[0].Interface()).Obj()
	}
	tuple := MakeTupleWithLen(len(results))
	for i := range results {
		tuple.Set(i, From(results[i].Interface()))
	}
	return tuple.Obj()
}

func FuncOf1(name string, fn unsafe.Pointer, doc string) Func {
	def := &C.PyMethodDef{
		ml_name:  AllocCStr(name),
		ml_meth:  C.PyCFunction(fn),
		ml_flags: C.METH_VARARGS,
		ml_doc:   AllocCStr(doc),
	}
	pyFn := C.PyCMethod_New(def, nil, nil, nil)
	return newFunc(pyFn)
}

func FuncOf(name string, fn any, doc string) Func {
	m := MainModule()
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		fmt.Printf("type: %T, kind: %d\n", fn, t.Kind())
		panic("AddFunction: fn must be a function")
	}
	ctx := &wrapperContext{v: fn, t: t}
	obj := C.PyCapsule_New(unsafe.Pointer(ctx), AllocCStr("wrapperContext"), nil)
	def := &C.PyMethodDef{
		ml_name:  AllocCStr(name),
		ml_meth:  C.PyCFunction(C.wrapperFunc),
		ml_flags: C.METH_VARARGS,
		ml_doc:   AllocCStr(doc),
	}
	pyFn := C.PyCMethod_New(def, obj, m.obj, nil)
	if pyFn == nil {
		panic(fmt.Sprintf("Failed to add function %s to module", name))
	}
	return newFunc(pyFn)
}
