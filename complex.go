package gp

/*
#include <Python.h>
*/
import "C"

type Complex struct {
	Object
}

func newComplex(obj *PyObject) Complex {
	return Complex{newObject(obj)}
}

func MakeComplex(f complex128) Complex {
	return newComplex(C.PyComplex_FromDoubles(C.double(real(f)), C.double(imag(f))))
}

func (c Complex) Complex128() complex128 {
	real := C.PyComplex_RealAsDouble(c.obj)
	imag := C.PyComplex_ImagAsDouble(c.obj)
	return complex(real, imag)
}

func (c Complex) Real() float64 {
	return float64(C.PyComplex_RealAsDouble(c.obj))
}

func (c Complex) Imag() float64 {
	return float64(C.PyComplex_ImagAsDouble(c.obj))
}
