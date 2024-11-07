package gp

/*
#include <Python.h>
*/
import "C"

import (
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

// ----------------------------------------------------------------------------

type holderList struct {
	head *objectHolder
}

func (l *holderList) PushFront(holder *objectHolder) {
	if l.head != nil {
		l.head.prev = holder
		holder.next = l.head
	}
	l.head = holder
}

func (l *holderList) Remove(holder *objectHolder) {
	if holder.prev != nil {
		holder.prev.next = holder.next
	} else {
		l.head = holder.next
	}
	if holder.next != nil {
		holder.next.prev = holder.prev
	}
}

// ----------------------------------------------------------------------------

const maxPyObjects = 128

type decRefList struct {
	objects []*C.PyObject
	mu      sync.Mutex
}

func (l *decRefList) add(obj *C.PyObject) {
	l.mu.Lock()
	l.objects = append(l.objects, obj)
	l.mu.Unlock()
}

func (l *decRefList) len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.objects)
}

func (l *decRefList) decRefAll() {
	l.mu.Lock()
	list := l.objects
	l.objects = make([]*C.PyObject, 0, maxPyObjects*2)
	l.mu.Unlock()

	for _, obj := range list {
		C.Py_DecRef(obj)
	}
}

// ----------------------------------------------------------------------------

type globalData struct {
	typeMetas  map[*C.PyObject]*typeMeta
	pyTypes    map[reflect.Type]*C.PyObject
	holders    holderList
	decRefList decRefList
	finished   int32
}

var (
	global *globalData
)

func getGlobalData() *globalData {
	return global
}

func (gd *globalData) addDecRef(obj *C.PyObject) {
	if atomic.LoadInt32(&gd.finished) != 0 {
		return
	}
	gd.decRefList.add(obj)
}

func (gd *globalData) decRefObjectsIfNeeded() {
	if gd.decRefList.len() >= maxPyObjects {
		gd.decRefList.decRefAll()
	}
}

// ----------------------------------------------------------------------------

func initGlobal() {
	global = &globalData{
		typeMetas: make(map[*C.PyObject]*typeMeta),
		pyTypes:   make(map[reflect.Type]*C.PyObject),
	}
}

func markFinished() {
	atomic.StoreInt32(&global.finished, 1)
}

func cleanupGlobal() {
	for _, meta := range global.typeMetas {
		for _, method := range meta.methods {
			def := method.def
			if def != nil {
				C.free(unsafe.Pointer(def.ml_name))
				C.free(unsafe.Pointer(def.ml_doc))
				C.free(unsafe.Pointer(def))
			}
		}
	}
	global = nil
}
