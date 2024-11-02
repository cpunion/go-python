package gp

/*
#include <Python.h>
#include <structmember.h>
#include <moduleobject.h>

#include "wrap.h"
extern PyObject* wrapperAlloc(PyTypeObject* type, Py_ssize_t size);
extern int wrapperInit(PyObject* self, PyObject* args);
static int isModule(PyObject* ob)
{
	return PyObject_TypeCheck(ob, &PyModule_Type);
}
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

func CreateFunc(fn any, doc string) Func {
	m := MainModule()
	return m.AddMethod("", fn, doc)
}

type wrapperType struct {
	PyObject
	v byte
}

type slotMetaType int

const (
	slotMethod slotMetaType = iota
	slotGet
	slotSet
)

type slotMeta struct {
	name       string
	methodName string
	fn         any
	doc        string
	hasRecv    bool         // whether it has a receiver
	slotType   slotMetaType // slot type
	index      int          // used for member type
	typ        reflect.Type // member/method type
}

type typeMeta struct {
	typ     reflect.Type
	wrapTyp reflect.Type
	init    *slotMeta
	methods map[uint]*slotMeta
}

var (
	typeMetaMap = make(map[*C.PyObject]*typeMeta)
	pyTypeMap   = make(map[reflect.Type]*C.PyObject)
)

//export wrapperAlloc
func wrapperAlloc(typ *C.PyTypeObject, size C.Py_ssize_t) *C.PyObject {
	self := C.PyType_GenericAlloc(typ, size)
	if self != nil {
		meta := typeMetaMap[(*C.PyObject)(unsafe.Pointer(typ))]
		wrapper := (*wrapperType)(unsafe.Pointer(self))
		vPtr := unsafe.Pointer(&wrapper.v)
		reflect.NewAt(meta.typ, vPtr).Elem().Set(reflect.New(meta.typ).Elem())
	}
	return self
}

//export wrapperInit
func wrapperInit(self, args *C.PyObject) C.int {
	typ := (*C.PyObject)(self).ob_type
	typeMeta := typeMetaMap[(*C.PyObject)(unsafe.Pointer(typ))]
	if typeMeta.init == nil {
		return 0
	}
	if wrapperMethod_(typeMeta, typeMeta.init, self, args, 0) == nil {
		return -1
	}
	return 0
}

//export getterMethod
func getterMethod(self *C.PyObject, _closure unsafe.Pointer, methodId C.int) *C.PyObject {
	typeMeta := typeMetaMap[(*C.PyObject)(unsafe.Pointer(self.ob_type))]
	if typeMeta == nil {
		SetError(fmt.Errorf("type %v not registered", FromPy(self)))
		return nil
	}
	methodMeta := typeMeta.methods[uint(methodId)]
	if methodMeta == nil {
		SetError(fmt.Errorf("getter method %d not found", methodId))
		return nil
	}
	if methodMeta.slotType != slotGet {
		SetError(fmt.Errorf("method %d is not a getter method", methodId))
		return nil
	}
	wrapper := (*wrapperType)(unsafe.Pointer(self))
	vPtr := unsafe.Pointer(&wrapper.v)
	goValue := reflect.NewAt(typeMeta.typ, vPtr).Elem()
	field := goValue.Field(methodMeta.index)
	return From(field.Interface()).Obj()
}

//export setterMethod
func setterMethod(self, value *C.PyObject, _closure unsafe.Pointer, methodId C.int) C.int {
	typeMeta := typeMetaMap[(*C.PyObject)(unsafe.Pointer(self.ob_type))]
	if typeMeta == nil {
		SetError(fmt.Errorf("type %v not registered", FromPy(self)))
		return -1
	}
	methodMeta := typeMeta.methods[uint(methodId)]
	if methodMeta == nil {
		SetError(fmt.Errorf("setter method %d not found", methodId))
		return -1
	}
	if methodMeta.slotType != slotSet {
		SetError(fmt.Errorf("method %d is not a setter method", methodId))
		return -1
	}
	wrapper := (*wrapperType)(unsafe.Pointer(self))
	vPtr := unsafe.Pointer(&wrapper.v)
	goValue := reflect.NewAt(typeMeta.typ, vPtr).Elem()
	field := goValue.Field(methodMeta.index)
	if !ToValue(FromPy(value), field) {
		SetError(fmt.Errorf("failed to convert value to %s", methodMeta.typ))
		return -1
	}
	return 0
}

//export wrapperMethod
func wrapperMethod(self, args *C.PyObject, methodId C.int) *C.PyObject {
	key := self
	if C.isModule(self) == 0 {
		key = (*C.PyObject)(unsafe.Pointer(self.ob_type))
	}

	typeMeta, ok := typeMetaMap[key]
	if !ok {
		SetError(fmt.Errorf("type %v not registered", FromPy(key)))
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

	methodType := reflect.TypeOf(methodMeta.fn)

	argc := C.PyTuple_Size(args)
	expectedArgs := methodType.NumIn()
	hasReceiver := methodMeta.hasRecv

	if hasReceiver {
		expectedArgs-- // decrease expected number if it has a receiver
	}

	if int(argc) != expectedArgs {
		SetTypeError(fmt.Errorf("method %s expects %d arguments, got %d", methodMeta.name, expectedArgs, argc))
		return nil
	}

	goArgs := make([]reflect.Value, methodType.NumIn())
	argIndex := 0

	if hasReceiver {
		wrapper := (*wrapperType)(unsafe.Pointer(self))
		vPtr := unsafe.Pointer(&wrapper.v)
		recv := reflect.NewAt(typeMeta.typ, vPtr)

		if methodType.In(0).Kind() == reflect.Struct {
			// Method expects value receiver
			recv = recv.Elem()
		}
		goArgs[0] = recv
		argIndex = 1
	}

	for i := 0; i < int(argc); i++ {
		arg := C.PyTuple_GetItem(args, C.Py_ssize_t(i))
		C.Py_IncRef(arg)
		argType := methodType.In(i + argIndex)
		argPy := FromPy(arg)
		goValue := reflect.New(argType).Elem()
		if !ToValue(argPy, goValue) {
			SetTypeError(fmt.Errorf("failed to convert argument %v to %v", argPy, argType))
			return nil
		}
		goArgs[i+argIndex] = goValue
	}

	results := reflect.ValueOf(methodMeta.fn).Call(goArgs)

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
				typ:        method.Type,
				hasRecv:    true,
				slotType:   slotMethod,
			}

			methodPtr := C.wrapperMethods[methodId]

			ret = append(ret, C.PyMethodDef{
				ml_name:  AllocCStrDontFree(pythonName),
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

func getMembers(t reflect.Type, methods map[uint]*slotMeta) (members *C.PyMemberDef, getsets *C.PyGetSetDef) {
	baseOffset := unsafe.Offsetof(wrapperType{}.v)
	membersList := make([]C.PyMemberDef, 0)
	getsetsList := make([]C.PyGetSetDef, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		pythonName := goNameToPythonName(field.Name)
		memberType := getMemberType(field.Type)

		if memberType != -1 {
			// create as member variable for C-compatible types
			membersList = append(membersList, C.PyMemberDef{
				name:   AllocCStrDontFree(pythonName),
				_type:  memberType,
				offset: C.Py_ssize_t(baseOffset + field.Offset),
			})
		} else {
			// if _, ok := pyTypeMap[field.Type]; !ok {
			// 	AddType(field.Type, nil, "", "")
			// }
			getId := uint(len(methods))
			methods[getId] = &slotMeta{
				name:       field.Name,
				methodName: pythonName,
				typ:        field.Type,
				slotType:   slotGet,
				hasRecv:    false,
				index:      i,
			}
			setId := uint(len(methods))
			methods[setId] = &slotMeta{
				name:       field.Name,
				methodName: pythonName,
				typ:        field.Type,
				slotType:   slotSet,
				hasRecv:    false,
				index:      i,
			}
			getsetsList = append(getsetsList, C.PyGetSetDef{
				name:    AllocCStrDontFree(pythonName),
				get:     C.getterMethods[getId],
				set:     C.setterMethods[setId],
				doc:     nil,
				closure: nil,
			})
		}
	}

	// Add null terminators
	membersList = append(membersList, C.PyMemberDef{})
	getsetsList = append(getsetsList, C.PyGetSetDef{})

	// Allocate and copy members array
	memberSize := C.size_t(C.sizeof_PyMemberDef * len(membersList))
	membersPtr := (*C.PyMemberDef)(C.malloc(memberSize))
	C.memset(unsafe.Pointer(membersPtr), 0, memberSize)

	memberArrayPtr := unsafe.Pointer(membersPtr)
	for i, member := range membersList {
		currentMember := (*C.PyMemberDef)(unsafe.Pointer(uintptr(memberArrayPtr) + uintptr(i)*unsafe.Sizeof(C.PyMemberDef{})))
		*currentMember = member
	}

	// Allocate and copy getsets array
	getsetSize := C.size_t(C.sizeof_PyGetSetDef * len(getsetsList))
	getsetsPtr := (*C.PyGetSetDef)(C.malloc(getsetSize))
	C.memset(unsafe.Pointer(getsetsPtr), 0, getsetSize)

	getsetArrayPtr := unsafe.Pointer(getsetsPtr)
	for i, getset := range getsetsList {
		currentGetSet := (*C.PyGetSetDef)(unsafe.Pointer(uintptr(getsetArrayPtr) + uintptr(i)*unsafe.Sizeof(C.PyGetSetDef{})))
		*currentGetSet = getset
	}

	return membersPtr, getsetsPtr
}

func (m Module) AddType(obj, init any, name, doc string) Object {
	ty := reflect.TypeOf(obj)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	if ty.Kind() != reflect.Struct {
		panic("AddType: obj must be a struct or pointer to struct")
	}

	wrapper := wrapperType{}
	meta := &typeMeta{
		typ:     ty,
		wrapTyp: reflect.TypeOf(wrapper),
		methods: make(map[uint]*slotMeta),
	}

	slots := make([]C.PyType_Slot, 0)
	if init != nil {
		slots = append(slots, C.PyType_Slot{slot: C.Py_tp_init, pfunc: unsafe.Pointer(C.wrapperInit)})
		meta.init = &slotMeta{
			name:       runtime.FuncForPC(reflect.ValueOf(init).Pointer()).Name(),
			methodName: "__init__",
			fn:         init,
			typ:        reflect.TypeOf(init),
			slotType:   slotMethod,
			hasRecv:    true,
		}
	}
	members, getsets := getMembers(ty, meta.methods)
	slots = append(slots, C.PyType_Slot{slot: C.Py_tp_members, pfunc: unsafe.Pointer(members)})
	slots = append(slots, C.PyType_Slot{slot: C.Py_tp_getset, pfunc: unsafe.Pointer(getsets)})
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

	totalSize := unsafe.Offsetof(wrapper.v) + ty.Size()
	spec := &C.PyType_Spec{
		name:      C.CString(name),
		basicsize: C.int(totalSize),
		flags:     C.Py_TPFLAGS_DEFAULT,
		slots:     slotsPtr,
	}

	typeObj := C.PyType_FromSpec(spec)
	if typeObj == nil {
		panic(fmt.Sprintf("Failed to create type %s", name))
	}

	typeMetaMap[typeObj] = meta
	pyTypeMap[ty] = typeObj

	if C.PyModule_AddObject(m.obj, C.CString(name), typeObj) < 0 {
		C.Py_DecRef(typeObj)
		panic(fmt.Sprintf("Failed to add type %s to module", name))
	}

	return newObject(typeObj)
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

	meta, ok := typeMetaMap[m.obj]
	if !ok {
		meta = &typeMeta{
			methods: make(map[uint]*slotMeta),
		}
		typeMetaMap[m.obj] = meta
	}

	methodId := uint(len(meta.methods))
	meta.methods[methodId] = &slotMeta{
		name:       name,
		methodName: name,
		fn:         fn,
		typ:        t,
		doc:        doc,
		slotType:   slotMethod,
		hasRecv:    false,
	}

	methodPtr := C.wrapperMethods[methodId]
	cName := C.CString(name)
	cDoc := C.CString(doc)

	def := &C.PyMethodDef{
		ml_name:  cName,
		ml_meth:  C.PyCFunction(methodPtr),
		ml_flags: C.METH_VARARGS,
		ml_doc:   cDoc,
	}

	pyFunc := C.PyCFunction_New(def, m.obj)
	if pyFunc == nil {
		panic(fmt.Sprintf("Failed to create function %s", name))
	}

	if C.PyModule_AddObject(m.obj, cName, pyFunc) < 0 {
		C.Py_DecRef(pyFunc)
		panic(fmt.Sprintf("Failed to add function %s to module", name))
	}

	return newFunc(pyFunc)
}

func SetError(err error) {
	errStr := C.CString(err.Error())
	C.PyErr_SetString(C.PyExc_RuntimeError, errStr)
	C.free(unsafe.Pointer(errStr))
}

func SetTypeError(err error) {
	errStr := C.CString(err.Error())
	C.PyErr_SetString(C.PyExc_TypeError, errStr)
	C.free(unsafe.Pointer(errStr))
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
