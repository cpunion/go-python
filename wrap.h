#ifndef __WRAP_H__
#define __WRAP_H__

#include <Python.h>

extern PyObject *wrapperMethod(PyObject *self, PyObject *args, int methodId);
extern PyObject *(*wrapperMethods[256])(PyObject *self, PyObject *args);

extern PyObject *getterMethod(PyObject *self, void *closure, int methodId);
extern int setterMethod(PyObject *self, PyObject *value, void *closure, int methodId);

extern getter getterMethods[256];
extern setter setterMethods[256];

#define DECLARE_GETTER_METHOD(ida, idb) \
    extern PyObject* getterMethod##ida##idb(PyObject* self, void* closure);

#define DECLARE_SETTER_METHOD(ida, idb) \
    extern int setterMethod##ida##idb(PyObject* self, PyObject* value, void* closure);

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


#define DECLARE_GETTER_METHODS(ida) \
    DECLARE_GETTER_METHOD(ida, 0) \
    DECLARE_GETTER_METHOD(ida, 1) \
    DECLARE_GETTER_METHOD(ida, 2) \
    DECLARE_GETTER_METHOD(ida, 3) \
    DECLARE_GETTER_METHOD(ida, 4) \
    DECLARE_GETTER_METHOD(ida, 5) \
    DECLARE_GETTER_METHOD(ida, 6) \
    DECLARE_GETTER_METHOD(ida, 7) \
    DECLARE_GETTER_METHOD(ida, 8) \
    DECLARE_GETTER_METHOD(ida, 9) \
    DECLARE_GETTER_METHOD(ida, a) \
    DECLARE_GETTER_METHOD(ida, b) \
    DECLARE_GETTER_METHOD(ida, c) \
    DECLARE_GETTER_METHOD(ida, d) \
    DECLARE_GETTER_METHOD(ida, e) \
    DECLARE_GETTER_METHOD(ida, f)

#define DECLARE_SETTER_METHODS(ida) \
    DECLARE_SETTER_METHOD(ida, 0) \
    DECLARE_SETTER_METHOD(ida, 1) \
    DECLARE_SETTER_METHOD(ida, 2) \
    DECLARE_SETTER_METHOD(ida, 3) \
    DECLARE_SETTER_METHOD(ida, 4) \
    DECLARE_SETTER_METHOD(ida, 5) \
    DECLARE_SETTER_METHOD(ida, 6) \
    DECLARE_SETTER_METHOD(ida, 7) \
    DECLARE_SETTER_METHOD(ida, 8) \
    DECLARE_SETTER_METHOD(ida, 9) \
    DECLARE_SETTER_METHOD(ida, a) \
    DECLARE_SETTER_METHOD(ida, b) \
    DECLARE_SETTER_METHOD(ida, c) \
    DECLARE_SETTER_METHOD(ida, d) \
    DECLARE_SETTER_METHOD(ida, e) \
    DECLARE_SETTER_METHOD(ida, f)

#define DECLARE_WRAPPER_ALL_GETTERS() \
    DECLARE_GETTER_METHODS(0) \
    DECLARE_GETTER_METHODS(1) \
    DECLARE_GETTER_METHODS(2) \
    DECLARE_GETTER_METHODS(3) \
    DECLARE_GETTER_METHODS(4) \
    DECLARE_GETTER_METHODS(5) \
    DECLARE_GETTER_METHODS(6) \
    DECLARE_GETTER_METHODS(7) \
    DECLARE_GETTER_METHODS(8) \
    DECLARE_GETTER_METHODS(9) \
    DECLARE_GETTER_METHODS(a) \
    DECLARE_GETTER_METHODS(b) \
    DECLARE_GETTER_METHODS(c) \
    DECLARE_GETTER_METHODS(d) \
    DECLARE_GETTER_METHODS(e) \
    DECLARE_GETTER_METHODS(f)

#define DECLARE_WRAPPER_ALL_SETTERS() \
    DECLARE_SETTER_METHODS(0) \
    DECLARE_SETTER_METHODS(1) \
    DECLARE_SETTER_METHODS(2) \
    DECLARE_SETTER_METHODS(3) \
    DECLARE_SETTER_METHODS(4) \
    DECLARE_SETTER_METHODS(5) \
    DECLARE_SETTER_METHODS(6) \
    DECLARE_SETTER_METHODS(7) \
    DECLARE_SETTER_METHODS(8) \
    DECLARE_SETTER_METHODS(9) \
    DECLARE_SETTER_METHODS(a) \
    DECLARE_SETTER_METHODS(b) \
    DECLARE_SETTER_METHODS(c) \
    DECLARE_SETTER_METHODS(d) \
    DECLARE_SETTER_METHODS(e) \
    DECLARE_SETTER_METHODS(f)

DECLARE_WRAPPER_ALL_GETTERS()
DECLARE_WRAPPER_ALL_SETTERS()

#endif
