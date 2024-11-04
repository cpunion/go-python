package main

import (
	"fmt"

	gp "github.com/cpunion/go-python"
	"github.com/cpunion/go-python/_demo/module/foo"
)

func main() {
	gp.Initialize()
	defer gp.Finalize()
	fooMod := foo.InitFooModule()
	gp.GetModuleDict().SetString("foo", fooMod)

	Main1(fooMod)
	Main2()
}

func Main1(fooMod gp.Module) {
	sum := fooMod.Call("add", 1, 2).AsLong()
	fmt.Printf("Sum of 1 + 2: %d\n", sum.Int64())

	dict := fooMod.Dict()
	Point := dict.Get(gp.MakeStr("Point")).AsFunc()

	point := Point.Call(3, 4)
	fmt.Printf("dir(point): %v\n", point.Dir())
	fmt.Printf("x: %v, y: %v\n", point.Attr("x"), point.Attr("y"))

	distance := point.Call("distance").AsFloat()
	fmt.Printf("Distance of 3 * 4: %f\n", distance.Float64())

	point.Call("move", 1, 2)
	fmt.Printf("x: %v, y: %v\n", point.Attr("x"), point.Attr("y"))

	distance = point.Call("distance").AsFloat()
	fmt.Printf("Distance of 4 * 6: %f\n", distance.Float64())

	point.Call("print")
}

func Main2() {
	fmt.Printf("=========== Main2 ==========\n")
	_ = gp.RunString(`
import foo
point = foo.Point(3, 4)
print("dir(point):", dir(point))
print("x:", point.x)
print("y:", point.y)

print("distance:", point.distance())

point.move(1, 2)
print("x:", point.x)
print("y:", point.y)
print("distance:", point.distance())

point.print()
	`)
}