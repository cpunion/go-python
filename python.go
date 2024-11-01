package gp

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type PyObject = C.PyObject
type PyCFunction = C.PyCFunction

func Initialize() {
	C.Py_Initialize()
}

func Finalize() {
	C.Py_FinalizeEx()
	typeMetaMap = make(map[*C.PyObject]*typeMeta)
	pyTypeMap = make(map[reflect.Type]*C.PyObject)
}

// ----------------------------------------------------------------------------

type InputType = C.int

const (
	SingleInput InputType = C.Py_single_input
	FileInput   InputType = C.Py_file_input
	EvalInput   InputType = C.Py_eval_input
)

func CompileString(code, filename string, start InputType) Object {
	return newObject(C.Py_CompileString(AllocCStr(code), AllocCStr(filename), C.int(start)))
}

func EvalCode(code Object, globals, locals Dict) Object {
	return newObject(C.PyEval_EvalCode(code.Obj(), globals.Obj(), locals.Obj()))
}

// ----------------------------------------------------------------------------

// llgo:link Cast llgo.staticCast
func Cast[U, T Objecter](obj T) (u U) {
	*(*T)(unsafe.Pointer(&u)) = obj
	return
}

// ----------------------------------------------------------------------------

func With[T Objecter](obj T, fn func(v T)) T {
	obj.object().Call("__enter__")
	defer obj.object().Call("__exit__")
	fn(obj)
	return obj
}

// ----------------------------------------------------------------------------

func MainModule() Module {
	return GetModule("__main__")
}

func None() Object {
	return newObject(C.Py_None)
}

func Nil() Object {
	return Object{}
}

// RunString executes Python code string and returns error if any
func RunString(code string) error {
	// Get __main__ module dict for executing code
	main := MainModule()
	dict := main.Dict()

	// Run the code string
	codeObj := CompileString(code, "<string>", FileInput)
	if codeObj.Nil() {
		return fmt.Errorf("failed to compile code")
	}

	ret := EvalCode(codeObj, dict, dict)
	if ret.Nil() {
		if err := FetchError(); err != nil {
			return err
		}
		return fmt.Errorf("failed to execute code")
	}
	return nil
}
