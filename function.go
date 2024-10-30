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

func newFunc(obj *C.PyObject) Func {
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
}

//export wrapperFunc
func wrapperFunc(self, args *C.PyObject) *C.PyObject {
	wCtx := (*wrapperContext)(C.PyCapsule_GetPointer(self, AllocCStr("wrapperContext")))
	fmt.Printf("wrapperContext: %p\n", wCtx)
	// 恢复上下文
	v := reflect.ValueOf(wCtx.v)
	t := v.Type()
	fmt.Printf("wrapperFunc type: %v\n", t)
	// 构建参数
	goArgs := make([]reflect.Value, t.NumIn())
	argsTuple := FromPy(args).AsTuple()
	fmt.Printf("args: %v\n", argsTuple)
	for i := range goArgs {
		goArgs[i] = reflect.New(t.In(i)).Elem()
		ToValue(FromPy(C.PyTuple_GetItem(args, C.Py_ssize_t(i))), goArgs[i])
		fmt.Printf("goArgs[%d]: %T\n", i, goArgs[i].Interface())
	}

	// 调用原始函数
	results := v.Call(goArgs)

	// 处理返回值
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

var ctxs = make(map[unsafe.Pointer]*wrapperContext)

func FuncOf(name string, fn any, doc string) Func {
	m := MainModule()
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		fmt.Printf("type: %T, kind: %d\n", fn, t.Kind())
		panic("AddFunction: fn must be a function")
	}
	println("FuncOf name:", name)
	fmt.Printf("FuncOf type: %v\n", t)
	ctx := new(wrapperContext)
	ctx.v = fn
	obj := C.PyCapsule_New(unsafe.Pointer(ctx), AllocCStr("wrapperContext"), nil)
	fmt.Printf("FuncOf ctx: %p\n", ctx)
	ctxs[unsafe.Pointer(ctx)] = ctx
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

func buildFormatString(t reflect.Type) *C.char {
	format := ""
	for i := 0; i < t.NumIn(); i++ {
		switch t.In(i).Kind() {
		case reflect.Int, reflect.Int64:
			format += "i"
		case reflect.Float64:
			format += "d"
		case reflect.String:
			format += "s"
		// Add more types as needed
		default:
			panic(fmt.Sprintf("Unsupported argument type: %v", t.In(i)))
		}
	}
	return AllocCStr(format)
}

func buildArgPointers(args []reflect.Value) []interface{} {
	pointers := make([]interface{}, len(args))
	for i := range args {
		args[i] = reflect.New(args[i].Type()).Elem()
		pointers[i] = args[i].Addr().Interface()
	}
	return pointers
}
