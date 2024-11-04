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

type threadData struct {
	typeMetas map[*C.PyObject]*typeMeta
	pyTypes   map[reflect.Type]*C.PyObject
}

var (
	globalThreadData sync.Map // map[int64]*threadData
)

func getCurrentThreadData() *threadData {
	id := getThreadID()
	maps, ok := globalThreadData.Load(id)
	if !ok {
		// 如果不存在，创建新的映射
		maps = &threadData{
			typeMetas: make(map[*C.PyObject]*typeMeta),
			pyTypes:   make(map[reflect.Type]*C.PyObject),
		}
		globalThreadData.Store(id, maps)
	}
	return maps.(*threadData)
}

// initThreadLocal 在线程初始化时调用
func initThreadLocal() {
	id := getThreadID()
	maps := &threadData{
		typeMetas: make(map[*C.PyObject]*typeMeta),
		pyTypes:   make(map[reflect.Type]*C.PyObject),
	}
	globalThreadData.Store(id, maps)
}

// cleanupThreadLocal 在线程结束时调用
func cleanupThreadLocal() {
	id := getThreadID()
	globalThreadData.Delete(id)
}

// getThreadID 获取当前线程ID
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
