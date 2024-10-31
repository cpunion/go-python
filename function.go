package gp

/*
#include <Python.h>
#include <structmember.h>

extern PyObject* wrapperFunc(PyObject* self, PyObject* args);

extern PyObject* wrapperMethod(PyObject* self, PyObject* args);
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unsafe"
)

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

func CreateFunc(fn any, doc string) Func {
	m := MainModule()
	return m.AddMethod("", fn, doc)
}

type wrapperType[T any] struct {
	PyObject
	v T
}

func newWrapper[T any](v T) wrapperType[T] {
	return wrapperType[T]{v: v}
}

type slotMeta struct {
	name string
	fn   any
	doc  string
}

type typeMeta struct {
	methods map[string]*slotMeta
}

var (
	typeMetaMap = make(map[*C.PyTypeObject]*typeMeta)
)

//export wrapperMethod
func wrapperMethod(self *C.PyObject, args *C.PyObject) *C.PyObject {
	return nil
}

func goNameToPythonName(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func getMethods_(t reflect.Type, methods map[string]*slotMeta) (ret []C.PyMethodDef) {
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.IsExported() {
			pythonName := goNameToPythonName(method.Name)
			methods[pythonName] = &slotMeta{fn: method.Func.Interface()}
			ret = append(ret, C.PyMethodDef{
				ml_name:  C.CString(pythonName),
				ml_meth:  (C.PyCFunction)(unsafe.Pointer(C.wrapperMethod)),
				ml_flags: C.METH_VARARGS,
			})
		}
	}
	return
}

func getMethods(t reflect.Type, methods map[string]*slotMeta) *C.PyMethodDef {
	methodsDef := getMethods_(t, methods)
	methodsDef = append(methodsDef, getMethods_(reflect.PointerTo(t), methods)...)
	methodsDef = append(methodsDef, C.PyMethodDef{})
	methodSize := C.size_t(C.sizeof_PyMethodDef * len(methodsDef))
	methodsPtr := (*C.PyMethodDef)(C.malloc(methodSize))
	C.memset(unsafe.Pointer(methodsPtr), 0, methodSize)

	methodArrayPtr := unsafe.Pointer(methodsPtr)
	for i, method := range methodsDef {
		currentMethod := (*C.PyMethodDef)(unsafe.Pointer(uintptr(methodArrayPtr) + uintptr(i)*unsafe.Sizeof(C.PyMethodDef{})))
		*currentMethod = method
	}
	return methodsPtr
}

func getMembers(t reflect.Type) (ret *C.PyMemberDef) {
	members := make([]C.PyMemberDef, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			pythonName := goNameToPythonName(field.Name)
			members = append(members, C.PyMemberDef{
				name:   C.CString(pythonName),
				_type:  C.T_INT,
				offset: C.Py_ssize_t(field.Offset),
			})
		}
	}
	members = append(members, C.PyMemberDef{})
	memberSize := C.size_t(C.sizeof_PyMemberDef * len(members))
	membersPtr := (*C.PyMemberDef)(C.malloc(memberSize))
	C.memset(unsafe.Pointer(membersPtr), 0, memberSize)

	memberArrayPtr := unsafe.Pointer(membersPtr)
	for i, member := range members {
		currentMember := (*C.PyMemberDef)(unsafe.Pointer(uintptr(memberArrayPtr) + uintptr(i)*unsafe.Sizeof(C.PyMemberDef{})))
		*currentMember = member
	}
	return membersPtr
}

func AddType[T any](m Module, init any, name string, doc string) Object {
	wrapper := wrapperType[T]{}
	ty := reflect.TypeOf(wrapper.v)
	if ty.Kind() != reflect.Struct {
		panic("AddType: t must be a struct")
	}

	meta := &typeMeta{
		methods: make(map[string]*slotMeta),
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	slots := make([]C.PyType_Slot, 0)
	slots = append(slots, C.PyType_Slot{slot: C.Py_tp_init, pfunc: C.wrapperMethod})
	meta.methods["__init__"] = &slotMeta{fn: init}
	slots = append(slots, C.PyType_Slot{slot: C.Py_tp_members, pfunc: unsafe.Pointer(getMembers(ty))})
	slots = append(slots, C.PyType_Slot{slot: C.Py_tp_methods, pfunc: unsafe.Pointer(getMethods(ty, meta.methods))})

	slotCount := len(slots) + 1
	slotSize := C.size_t(C.sizeof_PyType_Slot * slotCount)
	slotsPtr := (*C.PyType_Slot)(C.malloc(slotSize))
	C.memset(unsafe.Pointer(slotsPtr), 0, slotSize)

	slotArrayPtr := unsafe.Pointer(slotsPtr)
	for i, slot := range slots {
		currentSlot := (*C.PyType_Slot)(unsafe.Pointer(uintptr(slotArrayPtr) + uintptr(i)*unsafe.Sizeof(C.PyType_Slot{})))
		*currentSlot = slot
	}

	spec := &C.PyType_Spec{
		name:      cname,
		basicsize: C.int(unsafe.Sizeof(wrapper)),
		flags:     C.Py_TPFLAGS_DEFAULT,
		slots:     slotsPtr,
	}

	obj := C.PyType_FromSpec(spec)
	if obj == nil {
		panic(fmt.Sprintf("Failed to create type %s", name))
	}

	if C.PyModule_AddObject(m.obj, cname, obj) < 0 {
		C.Py_DecRef(obj)
		panic(fmt.Sprintf("Failed to add type %s to module", name))
	}

	return newObject(obj)
}

func (m Module) AddMethod(name string, fn any, doc string) Func {
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		fmt.Printf("type: %T, kind: %d\n", fn, t.Kind())
		panic("AddFunction: fn must be a function")
	}

	if name == "" {
		name = runtime.FuncForPC(v.Pointer()).Name()
	}
	if name == "" {
		name = fmt.Sprintf("anonymous_func_%p", fn)
	} else {
		if idx := strings.LastIndex(name, "."); idx >= 0 {
			name = name[idx+1:]
		}
	}

	doc = name + doc
	fmt.Printf("AddMethod: %s, %s\n", name, doc)

	// Create the wrapper context
	ctx := &wrapperContext{v: fn, t: t}

	// Create the capsule
	capsule := C.PyCapsule_New(unsafe.Pointer(ctx), AllocCStr("wrapperContext"), nil)

	// Create method definition with C-allocated strings
	cName := C.CString(name)
	cDoc := C.CString(doc)
	def := &C.PyMethodDef{
		ml_name:  cName,
		ml_meth:  C.PyCFunction(C.wrapperFunc),
		ml_flags: C.METH_VARARGS,
		ml_doc:   cDoc,
	}

	// Create the Python method using PyCMethod_New
	pyFn := C.PyCMethod_New(def, capsule, m.obj, nil)
	if pyFn == nil {
		panic(fmt.Sprintf("Failed to create function %s", name))
	}

	// Add the function to the module
	if C.PyModule_AddObject(m.obj, cName, pyFn) < 0 {
		panic(fmt.Sprintf("Failed to add function %s to module", name))
	}

	return newFunc(pyFn)
}
