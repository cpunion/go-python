package gp

/*
#include <Python.h>
#include <moduleobject.h>

static PyModuleDef_Base moduleHeadInit() {
	PyModuleDef_Base base = PyModuleDef_HEAD_INIT;
	return base;
}

*/
import "C"
import (
	"fmt"
)

// ModuleBuilder helps to build Python modules
type ModuleBuilder struct {
	name    string
	doc     string
	methods []C.PyMethodDef
}

// NewModuleBuilder creates a new ModuleBuilder
func NewModuleBuilder(name, doc string) *ModuleBuilder {
	return &ModuleBuilder{
		name: name,
		doc:  doc,
	}
}

// /* Flag passed to newmethodobject */
// /* #define METH_OLDARGS  0x0000   -- unsupported now */
// #define METH_VARARGS  0x0001
// #define METH_KEYWORDS 0x0002
// /* METH_NOARGS and METH_O must not be combined with the flags above. */
// #define METH_NOARGS   0x0004
// #define METH_O        0x0008

// /* METH_CLASS and METH_STATIC are a little different; these control
//    the construction of methods for a class.  These cannot be used for
//    functions in modules. */
// #define METH_CLASS    0x0010
// #define METH_STATIC   0x0020

// /* METH_COEXIST allows a method to be entered even though a slot has
//    already filled the entry.  When defined, the flag allows a separate
//    method, "__contains__" for example, to coexist with a defined
//    slot like sq_contains. */

// #define METH_COEXIST   0x0040

// #if !defined(Py_LIMITED_API) || Py_LIMITED_API+0 >= 0x030a0000
// #  define METH_FASTCALL  0x0080
// #endif

// /* This bit is preserved for Stackless Python */
// #ifdef STACKLESS
// #  define METH_STACKLESS 0x0100
// #else
// #  define METH_STACKLESS 0x0000
// #endif

// /* METH_METHOD means the function stores an
//  * additional reference to the class that defines it;
//  * both self and class are passed to it.
//  * It uses PyCMethodObject instead of PyCFunctionObject.
//  * May not be combined with METH_NOARGS, METH_O, METH_CLASS or METH_STATIC.
//  */

// #if !defined(Py_LIMITED_API) || Py_LIMITED_API+0 >= 0x03090000
// #define METH_METHOD 0x0200
// #endif

const (
	METH_VARARGS  = 0x0001
	METH_KEYWORDS = 0x0002
	METH_NOARGS   = 0x0004
	METH_O        = 0x0008
	METH_CLASS    = 0x0010
	METH_STATIC   = 0x0020
	METH_COEXIST  = 0x0040
	METH_FASTCALL = 0x0080
	METH_METHOD   = 0x0200
)

// AddMethod adds a method to the module
func (mb *ModuleBuilder) AddMethod(name string, fn PyCFunction, doc string) *ModuleBuilder {
	mb.methods = append(mb.methods, C.PyMethodDef{
		ml_name:  AllocCStr(name),
		ml_meth:  fn,
		ml_flags: METH_VARARGS,
		ml_doc:   AllocCStr(doc),
	})
	return mb
}

// Build creates and returns a new Python module
func (mb *ModuleBuilder) Build() Module {
	// Add a null terminator to the methods slice
	mb.methods = append(mb.methods, C.PyMethodDef{})
	def := &C.PyModuleDef{
		m_base:    C.moduleHeadInit(),
		m_name:    AllocCStr(mb.name),
		m_doc:     AllocCStr(mb.doc),
		m_size:    -1,
		m_methods: &mb.methods[0],
	}
	fmt.Printf("name: %s, doc: %s, size: %d\n", GoString(def.m_name), GoString(def.m_doc), def.m_size)
	for _, m := range mb.methods {
		fmt.Printf("method: %s, doc: %s\n", GoString(m.ml_name), GoString(m.ml_doc))
	}

	m := C.PyModule_Create2(def, 1013)

	if m == nil {
		panic("failed to create module")
	}

	return newModule(m)
}
