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
