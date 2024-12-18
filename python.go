package gp

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type cPyObject = C.PyObject

func Initialize() {
	runtime.LockOSThread()
	C.Py_Initialize()
	initGlobal()
}

func Finalize() {
	markFinished()
	r := C.Py_FinalizeEx()
	cleanupGlobal()
	check(r == 0, "failed to finalize Python")
}

// ----------------------------------------------------------------------------

type InputType = C.int

const (
	SingleInput InputType = C.Py_single_input
	FileInput   InputType = C.Py_file_input
	EvalInput   InputType = C.Py_eval_input
)

func CompileString(code, filename string, start InputType) (Object, error) {
	ccode := AllocCStr(code)
	cfilename := AllocCStr(filename)
	o := C.Py_CompileString(ccode, cfilename, C.int(start))
	// TODO: check why double free
	C.free(unsafe.Pointer(ccode))
	C.free(unsafe.Pointer(cfilename))
	if o == nil {
		err := FetchError()
		if err != nil {
			return Object{}, err
		}
		return Object{}, fmt.Errorf("failed to compile code")
	}
	return newObject(o), nil
}

func EvalCode(code Object, globals, locals Dict) Object {
	return newObject(C.PyEval_EvalCode(code.cpyObj(), globals.cpyObj(), locals.cpyObj()))
}

// ----------------------------------------------------------------------------

func cast[U, T Objecter](from T) (to U) {
	*(*T)(unsafe.Pointer(&to)) = from
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
	return ImportModule("__main__")
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
	check(!main.Nil(), "failed to get __main__ module")
	dict := main.Dict()

	// Run the code string
	codeObj, err := CompileString(code, "<string>", FileInput)
	if err != nil {
		return err
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

func RunMain(args []string) int {
	argc := len(args)
	argv := make([]*C.char, argc+1)
	for i, arg := range args {
		argv[i] = AllocCStr(arg)
	}
	return int(C.Py_BytesMain(C.int(argc), &argv[0]))
}

// ----------------------------------------------------------------------------

func check(b bool, msg string) {
	if !b {
		panic(msg)
	}
}
