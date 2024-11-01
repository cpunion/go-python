package gp

/*
#include <Python.h>
#include <structmember.h>

#include "wrap.h"
extern PyObject* wrapperFunc(PyObject* self, PyObject* args);
extern PyObject* wrapperAlloc(PyTypeObject* type, Py_ssize_t size);
extern int wrapperInit(PyObject* self, PyObject* args);
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
	name       string
	methodName string
	fn         any
	doc        string
}

type typeMeta struct {
	typ     reflect.Type
	init    *slotMeta
	methods map[uint]*slotMeta
}

var (
	typeMetaMap = make(map[*C.PyTypeObject]*typeMeta)
)

//export wrapperAlloc
func wrapperAlloc(typ *C.PyTypeObject, size C.Py_ssize_t) *C.PyObject {
	self := C.PyType_GenericAlloc(typ, size)
	if self != nil {
		meta := typeMetaMap[typ]
		wrapper := (*wrapperType[any])(unsafe.Pointer(self))

		wrapperVal := reflect.ValueOf(&wrapper.v).Elem()
		wrapperVal.Set(reflect.New(meta.typ).Elem())
	}
	return self
}

//export wrapperInit
func wrapperInit(self, args *C.PyObject) C.int {
	typ := (*C.PyObject)(self).ob_type
	typeObj := (*C.PyTypeObject)(unsafe.Pointer(typ))
	typeMeta := typeMetaMap[typeObj]
	if typeMeta.init == nil {
		return 0
	}
	if wrapperMethod_(typeMeta, typeMeta.init, self, args, 0) == nil {
		return -1
	}
	return 0
}

//export wrapperMethod
func wrapperMethod(self, args *C.PyObject, methodId C.int) *C.PyObject {
	// Get type object and metadata
	typ := (*C.PyObject)(self).ob_type
	typeObj := (*C.PyTypeObject)(unsafe.Pointer(typ))

	typeMeta, ok := typeMetaMap[typeObj]
	if !ok {
		SetError(fmt.Errorf("type %v not registered", typeObj))
		return nil
	}

	methodMeta := typeMeta.methods[uint(methodId)]
	return wrapperMethod_(typeMeta, methodMeta, self, args, methodId)
}

func wrapperMethod_(typeMeta *typeMeta, methodMeta *slotMeta, self, args *C.PyObject, methodId C.int) *C.PyObject {
	if methodMeta == nil {
		SetError(fmt.Errorf("method %d not found", methodId))
		return nil
	}

	// Get the wrapper and method type
	wrapper := (*wrapperType[any])(unsafe.Pointer(self))
	methodType := reflect.TypeOf(methodMeta.fn)

	// Get the address of wrapper.v and create reflect.Value
	vPtr := unsafe.Pointer(&wrapper.v)
	// Create receiver value based on method type
	recv := reflect.NewAt(typeMeta.typ, vPtr)

	if methodType.In(0).Kind() == reflect.Struct {
		// Method expects value receiver
		recv = recv.Elem()
	}

	// Parse arguments
	argc := C.PyTuple_Size(args)
	if int(argc)+1 != methodType.NumIn() {
		SetTypeError(fmt.Errorf("method %s expects %d arguments, got %d", methodMeta.name, methodType.NumIn()-1, argc))
		return nil
	}
	goArgs := make([]reflect.Value, argc+1)
	goArgs[0] = recv
	for i := 0; i < int(argc); i++ {
		arg := C.PyTuple_GetItem(args, C.Py_ssize_t(i))
		goValue := reflect.New(methodType.In(i + 1)).Elem()
		ToValue(FromPy(arg), goValue)
		goArgs[i+1] = goValue
	}
	// Call the method with correct receiver
	results := reflect.ValueOf(methodMeta.fn).Call(goArgs)

	// Handle return values
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

func getMethods_(t reflect.Type, methods map[uint]*slotMeta) (ret []C.PyMethodDef) {
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.IsExported() {
			methodId := uint(len(methods))

			pythonName := goNameToPythonName(method.Name)
			methods[methodId] = &slotMeta{
				name:       method.Name,
				methodName: pythonName,
				fn:         method.Func.Interface(),
			}

			methodPtr := C.wrapperMethods[methodId]

			ret = append(ret, C.PyMethodDef{
				ml_name:  C.CString(pythonName),
				ml_meth:  (C.PyCFunction)(unsafe.Pointer(methodPtr)),
				ml_flags: C.METH_VARARGS,
				ml_doc:   nil,
			})
		}
	}
	return
}

func getMethods(t reflect.Type, methods map[uint]*slotMeta) *C.PyMethodDef {
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

// getMemberType returns the C member type for Go types that are compatible with Python/C
// Returns -1 for incompatible types
func getMemberType(t reflect.Type) C.int {
	switch t.Kind() {
	case reflect.Bool:
		return C.T_BOOL
	case reflect.Int8:
		return C.T_BYTE
	case reflect.Int16:
		return C.T_SHORT
	case reflect.Int32:
		return C.T_INT
	case reflect.Int64:
		return C.T_LONG
	case reflect.Int:
		if unsafe.Sizeof(int(0)) == unsafe.Sizeof(int32(0)) {
			return C.T_INT
		} else {
			return C.T_LONG
		}
	case reflect.Uint8:
		return C.T_UBYTE
	case reflect.Uint16:
		return C.T_USHORT
	case reflect.Uint32:
		return C.T_UINT
	case reflect.Uint64:
		return C.T_ULONG
	case reflect.Uint:
		if unsafe.Sizeof(uint(0)) == unsafe.Sizeof(uint32(0)) {
			return C.T_INT
		} else {
			return C.T_LONG
		}
	case reflect.Float32:
		return C.T_FLOAT
	case reflect.Float64:
		return C.T_DOUBLE
	default:
		return -1
	}
}

func getMembers(t reflect.Type) (ret *C.PyMemberDef) {
	baseOffset := unsafe.Offsetof(wrapperType[any]{}.v)
	members := make([]C.PyMemberDef, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		memberType := getMemberType(field.Type)
		if memberType == -1 {
			// Skip non-C-compatible types - these will need getset handlers
			continue
		}

		pythonName := goNameToPythonName(field.Name)
		members = append(members, C.PyMemberDef{
			name:   C.CString(pythonName),
			_type:  memberType,
			offset: C.Py_ssize_t(baseOffset + field.Offset),
		})
	}

	// Add null terminator
	members = append(members, C.PyMemberDef{})

	// Allocate and copy members array
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
		typ:     ty,
		methods: make(map[uint]*slotMeta),
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	slots := make([]C.PyType_Slot, 0)
	if init != nil {
		slots = append(slots, C.PyType_Slot{slot: C.Py_tp_init, pfunc: unsafe.Pointer(C.wrapperInit)})
		meta.init = &slotMeta{
			name:       runtime.FuncForPC(reflect.ValueOf(init).Pointer()).Name(),
			methodName: "__init__",
			fn:         init,
		}
	}
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

	// 将类型对象和meta信息保存到映射中
	typeObj := (*C.PyTypeObject)(unsafe.Pointer(obj))
	typeMetaMap[typeObj] = meta
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

func SetError(err error) {
	errStr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(errStr))
	C.PyErr_SetString(C.PyExc_RuntimeError, errStr)
}

func SetTypeError(err error) {
	errStr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(errStr))
	C.PyErr_SetString(C.PyExc_TypeError, errStr)
}

// FetchError returns the current Python error as a Go error
func FetchError() error {
	var ptype, pvalue, ptraceback *C.PyObject
	C.PyErr_Fetch(&ptype, &pvalue, &ptraceback)
	if ptype == nil {
		return nil
	}
	defer C.Py_DecRef(ptype)
	if pvalue == nil {
		return fmt.Errorf("python error")
	}
	defer C.Py_DecRef(pvalue)
	if ptraceback != nil {
		defer C.Py_DecRef(ptraceback)
	}

	// Convert error to string
	pyStr := C.PyObject_Str(pvalue)
	if pyStr == nil {
		return fmt.Errorf("python error (failed to get error message)")
	}
	defer C.Py_DecRef(pyStr)

	// Get error message as Go string
	cstr := C.PyUnicode_AsUTF8(pyStr)
	if cstr == nil {
		return fmt.Errorf("python error (failed to decode error message)")
	}

	return fmt.Errorf("python error: %s", C.GoString(cstr))
}
