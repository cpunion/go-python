#ifndef __WRAP_H__
#define __WRAP_H__

#include <Python.h>

extern PyObject *wrapperMethod(PyObject *self, PyObject *args, int methodId);
extern PyObject *(*wrapperMethods[256])(PyObject *self, PyObject *args);

#define DECLARE_WRAP_METHOD(ida, idb) \
	extern PyObject* wrapperMethod##ida##idb(PyObject* self, PyObject* args);

#define DECLARE_WRAP_METHODS(ida)                                              \
  DECLARE_WRAP_METHOD(ida, 0)                                                  \
  DECLARE_WRAP_METHOD(ida, 1)                                                  \
  DECLARE_WRAP_METHOD(ida, 2)                                                  \
  DECLARE_WRAP_METHOD(ida, 3)                                                  \
  DECLARE_WRAP_METHOD(ida, 4)                                                  \
  DECLARE_WRAP_METHOD(ida, 5)                                                  \
  DECLARE_WRAP_METHOD(ida, 6)                                                  \
  DECLARE_WRAP_METHOD(ida, 7)                                                  \
  DECLARE_WRAP_METHOD(ida, 8)                                                  \
  DECLARE_WRAP_METHOD(ida, 9)                                                  \
  DECLARE_WRAP_METHOD(ida, a)                                                  \
  DECLARE_WRAP_METHOD(ida, b)                                                  \
  DECLARE_WRAP_METHOD(ida, c)                                                  \
  DECLARE_WRAP_METHOD(ida, d)                                                  \
  DECLARE_WRAP_METHOD(ida, e)                                                  \
  DECLARE_WRAP_METHOD(ida, f)

#define DECLARE_WRAPPER_ALL_METHODS() \
	DECLARE_WRAP_METHODS(0) \
	DECLARE_WRAP_METHODS(1) \
	DECLARE_WRAP_METHODS(2) \
	DECLARE_WRAP_METHODS(3) \
	DECLARE_WRAP_METHODS(4) \
	DECLARE_WRAP_METHODS(5) \
	DECLARE_WRAP_METHODS(6) \
	DECLARE_WRAP_METHODS(7) \
	DECLARE_WRAP_METHODS(8) \
	DECLARE_WRAP_METHODS(9) \
	DECLARE_WRAP_METHODS(a) \
	DECLARE_WRAP_METHODS(b) \
	DECLARE_WRAP_METHODS(c) \
	DECLARE_WRAP_METHODS(d) \
	DECLARE_WRAP_METHODS(e) \
	DECLARE_WRAP_METHODS(f)

DECLARE_WRAPPER_ALL_METHODS()

#endif
