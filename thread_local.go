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

type threadData struct {
	typeMetas map[*C.PyObject]*typeMeta
	pyTypes   map[reflect.Type]*C.PyObject
	holders   holderList
}

var (
	globalThreadData sync.Map // map[int64]*threadData
)

func getCurrentThreadData() *threadData {
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
