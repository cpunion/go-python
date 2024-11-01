package foo

/*
#cgo pkg-config: python3-embed
#include <Python.h>
*/
import "C"

import (
	"fmt"

	gp "github.com/cpunion/go-python"
)

type Point struct {
	X float64
	Y float64
}

func (p *Point) init(x, y float64) {
	p.X = x
	p.Y = y
}

func (p *Point) Print() {
	fmt.Printf("Point(%f, %f)\n", p.X, p.Y)
}

func (p *Point) Distance() float64 {
	return p.X * p.Y
}

// Move method for Point
func (p *Point) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

func Add(a, b int) int {
	return a + b
}

func InitFooModule() gp.Module {
	m := gp.CreateModule("foo")
	fmt.Printf("CreateModule: %v\n", m)

	// Add the function to the module
	f := m.AddMethod("add", Add, "(a, b) -> int\n--\n\nAdd two integers.")
	fmt.Printf("AddMethod: %v\n", f)

	// Add the type to the module
	gp.AddType[Point](m, (*Point).init, "Point", "Point objects")

	return m
}
