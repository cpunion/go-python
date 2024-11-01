package main

import (
	"fmt"
	"runtime"

	gp "github.com/cpunion/go-python"
	"github.com/cpunion/go-python/_demo/autoderef/foo"
	pymath "github.com/cpunion/go-python/math"
)

func main() {
	gp.Initialize()
	defer gp.Finalize()
	fooMod := foo.InitFooModule()
	gp.GetModuleDict().Set(gp.MakeStr("foo").Object, fooMod.Object)

	Main1(fooMod)
	Main2()
	Main3()
}

func Main1(fooMod gp.Module) {
	fmt.Printf("=========== Main1 ==========\n")
	sum := fooMod.Call("add", gp.MakeLong(1), gp.MakeLong(2)).AsLong()
	fmt.Printf("Sum of 1 + 2: %d\n", sum.Int64())

	dict := fooMod.Dict()
	Point := dict.Get(gp.MakeStr("Point")).AsFunc()

	point := Point.Call(gp.MakeLong(3), gp.MakeLong(4))
	fmt.Printf("dir(point): %v\n", point.Dir())
	fmt.Printf("x: %v, y: %v\n", point.GetAttr("x"), point.GetAttr("y"))

	distance := point.Call("distance").AsFloat()
	fmt.Printf("Distance of 3 * 4: %f\n", distance.Float64())

	point.Call("move", gp.MakeFloat(1), gp.MakeFloat(2))
	fmt.Printf("x: %v, y: %v\n", point.GetAttr("x"), point.GetAttr("y"))

	distance = point.Call("distance").AsFloat()
	fmt.Printf("Distance of 4 * 6: %f\n", distance.Float64())
	point.Call("print")
}

func Main2() {
	fmt.Printf("=========== Main2 ==========\n")
	gp.RunString(`
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

func Main3() {
	pythonCode := `
def allocate_memory():
    return bytearray(10 * 1024 * 1024)

def memory_allocation_test():
    memory_blocks = []
    for i in range(10):
        memory_blocks.append(allocate_memory())
    print('Memory allocation test completed.')
    return memory_blocks

for i in range(10):
    memory_allocation_test()
`

	mod := gp.ImportModule("__main__")
	gbl := mod.Dict()
	code := gp.CompileString(pythonCode, "<string>", gp.FileInput)
	_ = gp.EvalCode(code, gbl, gp.Nil().AsDict())
	for i := 0; i < 10; i++ {
		result := gp.EvalCode(code, gbl, gp.Nil().AsDict())
		if result.Nil() {
			fmt.Printf("Failed to execute Python code\n")
			return
		}
		fmt.Printf("Iteration %d in python\n", i+1)
	}

	memory_allocation_test := mod.GetFuncAttr("memory_allocation_test")

	for i := 0; i < 100; i++ {
		// 100MB every time
		memory_allocation_test.Call()
		fmt.Printf("Iteration %d in go\n", i+1)
		runtime.GC()
	}

	for i := 1; i <= 100000; i++ {
		println(i)
		f := gp.MakeFloat(float64(i))
		r := pymath.Sqrt(f)
		b := r.IsInteger()
		var _ bool = b.Bool()
		if i%10000 == 0 {
			fmt.Printf("Iteration %d in go\n", i)
		}
	}

	gp.Finalize()
	fmt.Printf("Done\n")
}
