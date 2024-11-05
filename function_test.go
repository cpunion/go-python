package gp

import (
	"testing"
)

func TestFuncCall(t *testing.T) {
	Initialize()
	defer Finalize()

	// Get the Python built-in len function
	builtins := ImportModule("builtins")
	lenFunc := builtins.AttrFunc("len")

	func() {
		// Test calling with no args (should fail, len requires an argument)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when calling len() with no arguments")
			}
		}()
		lenFunc.Call()
		t.Error("Should not reach this point")
	}()

	// Test calling with one arg
	list := MakeList(1, 2, 3)
	result := lenFunc.Call(list)
	if result.AsLong().Int64() != 3 {
		t.Errorf("Expected len([1,2,3]) to be 3, got %d", result.AsLong().Int64())
	}

	// Test str.format with keyword args
	str := MakeStr("Hello {name}!")
	formatFunc := str.AttrFunc("format")
	result = formatFunc.Call(KwArgs{"name": "World"})
	if result.AsStr().String() != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got '%s'", result.AsStr().String())
	}

	// Test calling with multiple keyword args
	str = MakeStr("{greeting} {name}!")
	formatFunc = str.AttrFunc("format")
	result = formatFunc.Call(KwArgs{
		"greeting": "Hi",
		"name":     "Python",
	})
	if result.AsStr().String() != "Hi Python!" {
		t.Errorf("Expected 'Hi Python!', got '%s'", result.AsStr().String())
	}

	// Test pow function with positional args only
	math := ImportModule("math")
	powFunc := math.AttrFunc("pow")
	result = powFunc.Call(2, 3)
	if result.AsFloat().Float64() != 8 { // 2^3 = 8
		t.Errorf("Expected pow(2, 3) to be 8, got %v", result)
	}
}
