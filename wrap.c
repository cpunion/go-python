#include "wrap.h"


#define WRAP_METHOD(ida, idb)                                                  \
PyObject *wrapperMethod##ida##idb(PyObject *self, PyObject *args) {          \
	return wrapperMethod(self, args, 0x##ida * 16 + 0x##idb); \
}

#define WRAP_METHODS(ida)                                                      \
  WRAP_METHOD(ida, 0)                                                          \
  WRAP_METHOD(ida, 1)                                                          \
  WRAP_METHOD(ida, 2)                                                          \
  WRAP_METHOD(ida, 3)                                                          \
  WRAP_METHOD(ida, 4)                                                          \
  WRAP_METHOD(ida, 5)                                                          \
  WRAP_METHOD(ida, 6)                                                          \
  WRAP_METHOD(ida, 7)                                                          \
  WRAP_METHOD(ida, 8)                                                          \
  WRAP_METHOD(ida, 9)                                                          \
  WRAP_METHOD(ida, a)                                                          \
  WRAP_METHOD(ida, b)                                                          \
  WRAP_METHOD(ida, c)                                                          \
  WRAP_METHOD(ida, d)                                                          \
  WRAP_METHOD(ida, e)                                                          \
  WRAP_METHOD(ida, f)

#define WRAP_METHOD_ALL() \
	WRAP_METHODS(0) \
	WRAP_METHODS(1) \
	WRAP_METHODS(2) \
	WRAP_METHODS(3) \
	WRAP_METHODS(4) \
	WRAP_METHODS(5) \
	WRAP_METHODS(6) \
	WRAP_METHODS(7) \
	WRAP_METHODS(8) \
	WRAP_METHODS(9) \
	WRAP_METHODS(a) \
	WRAP_METHODS(b) \
	WRAP_METHODS(c) \
	WRAP_METHODS(d) \
	WRAP_METHODS(e) \
	WRAP_METHODS(f)

WRAP_METHOD_ALL()

#define WARP_METHOD_NAME(ida, idb) wrapperMethod##ida##idb,

#define WARP_METHOD_NAMES(ida) \
	WARP_METHOD_NAME(ida, 0) \
	WARP_METHOD_NAME(ida, 1) \
	WARP_METHOD_NAME(ida, 2) \
	WARP_METHOD_NAME(ida, 3) \
	WARP_METHOD_NAME(ida, 4) \
	WARP_METHOD_NAME(ida, 5) \
	WARP_METHOD_NAME(ida, 6) \
	WARP_METHOD_NAME(ida, 7) \
	WARP_METHOD_NAME(ida, 8) \
	WARP_METHOD_NAME(ida, 9) \
	WARP_METHOD_NAME(ida, a) \
	WARP_METHOD_NAME(ida, b) \
	WARP_METHOD_NAME(ida, c) \
	WARP_METHOD_NAME(ida, d) \
	WARP_METHOD_NAME(ida, e) \
	WARP_METHOD_NAME(ida, f)

#define WARP_METHOD_NAMES_ALL() \
	WARP_METHOD_NAMES(0) \
	WARP_METHOD_NAMES(1) \
	WARP_METHOD_NAMES(2) \
	WARP_METHOD_NAMES(3) \
	WARP_METHOD_NAMES(4) \
	WARP_METHOD_NAMES(5) \
	WARP_METHOD_NAMES(6) \
	WARP_METHOD_NAMES(7) \
	WARP_METHOD_NAMES(8) \
	WARP_METHOD_NAMES(9) \
	WARP_METHOD_NAMES(a) \
	WARP_METHOD_NAMES(b) \
	WARP_METHOD_NAMES(c) \
	WARP_METHOD_NAMES(d) \
	WARP_METHOD_NAMES(e) \
	WARP_METHOD_NAMES(f)

PyObject* (*wrapperMethods[256])(PyObject *self, PyObject *args) = {
	WARP_METHOD_NAMES_ALL()
};

#define GETTER_METHOD(ida, idb) \
PyObject *getterMethod##ida##idb(PyObject *self, void *closure) { \
    return getterMethod(self, closure, 0x##ida * 16 + 0x##idb); \
}

#define SETTER_METHOD(ida, idb) \
int setterMethod##ida##idb(PyObject *self, PyObject *value, void *closure) { \
    return setterMethod(self, value, closure, 0x##ida * 16 + 0x##idb); \
}

#define GETTER_METHODS(ida) \
    GETTER_METHOD(ida, 0) \
    GETTER_METHOD(ida, 1) \
    GETTER_METHOD(ida, 2) \
    GETTER_METHOD(ida, 3) \
    GETTER_METHOD(ida, 4) \
    GETTER_METHOD(ida, 5) \
    GETTER_METHOD(ida, 6) \
    GETTER_METHOD(ida, 7) \
    GETTER_METHOD(ida, 8) \
    GETTER_METHOD(ida, 9) \
    GETTER_METHOD(ida, a) \
    GETTER_METHOD(ida, b) \
    GETTER_METHOD(ida, c) \
    GETTER_METHOD(ida, d) \
    GETTER_METHOD(ida, e) \
    GETTER_METHOD(ida, f)

#define SETTER_METHODS(ida) \
    SETTER_METHOD(ida, 0) \
    SETTER_METHOD(ida, 1) \
    SETTER_METHOD(ida, 2) \
    SETTER_METHOD(ida, 3) \
    SETTER_METHOD(ida, 4) \
    SETTER_METHOD(ida, 5) \
    SETTER_METHOD(ida, 6) \
    SETTER_METHOD(ida, 7) \
    SETTER_METHOD(ida, 8) \
    SETTER_METHOD(ida, 9) \
    SETTER_METHOD(ida, a) \
    SETTER_METHOD(ida, b) \
    SETTER_METHOD(ida, c) \
    SETTER_METHOD(ida, d) \
    SETTER_METHOD(ida, e) \
    SETTER_METHOD(ida, f)

#define GETTER_METHOD_ALL() \
    GETTER_METHODS(0) \
    GETTER_METHODS(1) \
    GETTER_METHODS(2) \
    GETTER_METHODS(3) \
    GETTER_METHODS(4) \
    GETTER_METHODS(5) \
    GETTER_METHODS(6) \
    GETTER_METHODS(7) \
    GETTER_METHODS(8) \
    GETTER_METHODS(9) \
    GETTER_METHODS(a) \
    GETTER_METHODS(b) \
    GETTER_METHODS(c) \
    GETTER_METHODS(d) \
    GETTER_METHODS(e) \
    GETTER_METHODS(f)

#define SETTER_METHOD_ALL() \
    SETTER_METHODS(0) \
    SETTER_METHODS(1) \
    SETTER_METHODS(2) \
    SETTER_METHODS(3) \
    SETTER_METHODS(4) \
    SETTER_METHODS(5) \
    SETTER_METHODS(6) \
    SETTER_METHODS(7) \
    SETTER_METHODS(8) \
    SETTER_METHODS(9) \
    SETTER_METHODS(a) \
    SETTER_METHODS(b) \
    SETTER_METHODS(c) \
    SETTER_METHODS(d) \
    SETTER_METHODS(e) \
    SETTER_METHODS(f)

GETTER_METHOD_ALL()
SETTER_METHOD_ALL()

#define WARP_GETTER_METHOD_NAME(ida, idb) getterMethod##ida##idb,
#define WARP_SETTER_METHOD_NAME(ida, idb) setterMethod##ida##idb,

#define WARP_GETTER_METHOD_NAMES(ida) \
	WARP_GETTER_METHOD_NAME(ida, 0) \
	WARP_GETTER_METHOD_NAME(ida, 1) \
	WARP_GETTER_METHOD_NAME(ida, 2) \
	WARP_GETTER_METHOD_NAME(ida, 3) \
	WARP_GETTER_METHOD_NAME(ida, 4) \
	WARP_GETTER_METHOD_NAME(ida, 5) \
	WARP_GETTER_METHOD_NAME(ida, 6) \
	WARP_GETTER_METHOD_NAME(ida, 7) \
	WARP_GETTER_METHOD_NAME(ida, 8) \
	WARP_GETTER_METHOD_NAME(ida, 9) \
	WARP_GETTER_METHOD_NAME(ida, a) \
	WARP_GETTER_METHOD_NAME(ida, b) \
	WARP_GETTER_METHOD_NAME(ida, c) \
	WARP_GETTER_METHOD_NAME(ida, d) \
	WARP_GETTER_METHOD_NAME(ida, e) \
	WARP_GETTER_METHOD_NAME(ida, f)

#define WARP_SETTER_METHOD_NAMES(ida) \
	WARP_SETTER_METHOD_NAME(ida, 0) \
	WARP_SETTER_METHOD_NAME(ida, 1) \
	WARP_SETTER_METHOD_NAME(ida, 2) \
	WARP_SETTER_METHOD_NAME(ida, 3) \
	WARP_SETTER_METHOD_NAME(ida, 4) \
	WARP_SETTER_METHOD_NAME(ida, 5) \
	WARP_SETTER_METHOD_NAME(ida, 6) \
	WARP_SETTER_METHOD_NAME(ida, 7) \
	WARP_SETTER_METHOD_NAME(ida, 8) \
	WARP_SETTER_METHOD_NAME(ida, 9) \
	WARP_SETTER_METHOD_NAME(ida, a) \
	WARP_SETTER_METHOD_NAME(ida, b) \
	WARP_SETTER_METHOD_NAME(ida, c) \
	WARP_SETTER_METHOD_NAME(ida, d) \
	WARP_SETTER_METHOD_NAME(ida, e) \
	WARP_SETTER_METHOD_NAME(ida, f)

#define WARP_GETTER_METHOD_NAMES_ALL()                                         \
  WARP_GETTER_METHOD_NAMES(0)                                                   \
  WARP_GETTER_METHOD_NAMES(1)                                                   \
  WARP_GETTER_METHOD_NAMES(2)                                                   \
  WARP_GETTER_METHOD_NAMES(3)                                                   \
  WARP_GETTER_METHOD_NAMES(4)                                                   \
  WARP_GETTER_METHOD_NAMES(5)                                                   \
  WARP_GETTER_METHOD_NAMES(6)                                                   \
  WARP_GETTER_METHOD_NAMES(7)                                                   \
  WARP_GETTER_METHOD_NAMES(8)                                                   \
  WARP_GETTER_METHOD_NAMES(9)                                                   \
  WARP_GETTER_METHOD_NAMES(a)                                                   \
  WARP_GETTER_METHOD_NAMES(b)                                                   \
  WARP_GETTER_METHOD_NAMES(c)                                                   \
  WARP_GETTER_METHOD_NAMES(d)                                                   \
  WARP_GETTER_METHOD_NAMES(e)                                                   \
  WARP_GETTER_METHOD_NAMES(f)

#define WARP_SETTER_METHOD_NAMES_ALL() \
	WARP_SETTER_METHOD_NAMES(0) \
	WARP_SETTER_METHOD_NAMES(1) \
	WARP_SETTER_METHOD_NAMES(2) \
	WARP_SETTER_METHOD_NAMES(3) \
	WARP_SETTER_METHOD_NAMES(4) \
	WARP_SETTER_METHOD_NAMES(5) \
	WARP_SETTER_METHOD_NAMES(6) \
	WARP_SETTER_METHOD_NAMES(7) \
	WARP_SETTER_METHOD_NAMES(8) \
	WARP_SETTER_METHOD_NAMES(9) \
	WARP_SETTER_METHOD_NAMES(a) \
	WARP_SETTER_METHOD_NAMES(b) \
	WARP_SETTER_METHOD_NAMES(c) \
	WARP_SETTER_METHOD_NAMES(d) \
	WARP_SETTER_METHOD_NAMES(e) \
	WARP_SETTER_METHOD_NAMES(f)

getter getterMethods[256] = {	
	WARP_GETTER_METHOD_NAMES_ALL()
};

setter setterMethods[256] = {
	WARP_SETTER_METHOD_NAMES_ALL()
};
