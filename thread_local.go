package gp

/*
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

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

type decRefList struct {
	objects []*C.PyObject
	mu      sync.Mutex
}

func (l *decRefList) add(obj *C.PyObject) {
	l.mu.Lock()
	l.objects = append(l.objects, obj)
	l.mu.Unlock()
}

func (l *decRefList) decRefAll() {
	var list []*C.PyObject

	l.mu.Lock()
	list = l.objects
	l.objects = make([]*C.PyObject, 0, maxPyObjects*2)
	l.mu.Unlock()

	for _, obj := range list {
		C.Py_DecRef(obj)
	}
}

type threadData struct {
	typeMetas  map[*C.PyObject]*typeMeta
	pyTypes    map[reflect.Type]*C.PyObject
	holders    holderList
	decRefList decRefList
}

const maxPyObjects = 128

func (td *threadData) addPyObject(obj *C.PyObject) {
	td.decRefList.add(obj)
}

func (td *threadData) decRefObjectsIfNeeded() {
	if len(td.decRefList.objects) > maxPyObjects {
		td.decRefList.decRefAll()
	}
}

var (
	globalThreadData sync.Map // map[int64]*threadData
)

func getCurrentThreadData() *threadData {
	id := getThreadID()
	return getThreadData(id)
}

func getThreadData(gid int64) *threadData {
	id := getThreadID()
	maps, ok := globalThreadData.Load(id)
	if !ok {
		// if not exists, create new thread data
		maps = &threadData{
			typeMetas: make(map[*C.PyObject]*typeMeta),
			pyTypes:   make(map[reflect.Type]*C.PyObject),
		}
		globalThreadData.Store(id, maps)
	}
	return maps.(*threadData)
}

func initThreadLocal() {
	id := getThreadID()
	maps := &threadData{
		typeMetas: make(map[*C.PyObject]*typeMeta),
		pyTypes:   make(map[reflect.Type]*C.PyObject),
	}
	globalThreadData.Store(id, maps)
}

func cleanupThreadLocal() {
	id := getThreadID()
	globalThreadData.Delete(id)
}

func getThreadID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	id := int64(0)
	_, err := fmt.Sscanf(string(buf[:n]), "goroutine %d ", &id)
	if err != nil {
		panic(err)
	}
	return id
}
