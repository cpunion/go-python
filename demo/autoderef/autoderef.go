package main

import (
	"fmt"
	"runtime"

	. "github.com/gotray/go-python"
	pymath "github.com/gotray/go-python/math"
)

func main() {
	Initialize()
	defer Finalize()

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

	mod := ImportModule("__main__")
	gbl := mod.Dict()
	code, err := CompileString(pythonCode, "<string>", FileInput)
	if err != nil {
		fmt.Printf("Failed to compile Python code: %v\n", err)
		return
	}
	_ = EvalCode(code, gbl, Nil().AsDict())
	for i := 0; i < 10; i++ {
		result := EvalCode(code, gbl, Nil().AsDict())
		if result.Nil() {
			fmt.Printf("Failed to execute Python code\n")
			return
		}
		fmt.Printf("Iteration %d in python\n", i+1)
	}

	memory_allocation_test := mod.AttrFunc("memory_allocation_test")

	for i := 0; i < 100; i++ {
		// 100MB every time
		memory_allocation_test.Call()
		fmt.Printf("Iteration %d in go\n", i+1)
		runtime.GC()
	}

	for i := 1; i <= 1000000; i++ {
		f := MakeFloat(float64(i))
		_ = pymath.Sqrt(f)
		if i%10000 == 0 {
			fmt.Printf("Iteration %d in go\n", i)
		}
	}

	fmt.Printf("Done\n")
}
