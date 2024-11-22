package foo

import (
	"fmt"

	. "github.com/gotray/go-python"
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

func InitFooModule() Module {
	m := CreateModule("foo")
	// Add the function to the module
	m.AddMethod("add", Add, "(a, b) -> float\n--\n\nAdd two integers.")
	// Add the type to the module
	m.AddType(Point{}, (*Point).init, "Point", "Point objects")
	return m
}
