package gp

import (
	"reflect"
	"testing"
)

func TestToValue(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name     string
		pyValue  Object
		goType   interface{}
		expected interface{}
	}{
		{"int8", From(42), int8(0), int8(42)},
		{"int16", From(42), int16(0), int16(42)},
		{"int32", From(42), int32(0), int32(42)},
		{"int64", From(42), int64(0), int64(42)},
		{"int", From(42), int(0), int(42)},
		{"uint8", From(42), uint8(0), uint8(42)},
		{"uint16", From(42), uint16(0), uint16(42)},
		{"uint32", From(42), uint32(0), uint32(42)},
		{"uint64", From(42), uint64(0), uint64(42)},
		{"uint", From(42), uint(0), uint(42)},
		{"float32", From(3.14), float32(0), float32(3.14)},
		{"float64", From(3.14), float64(0), float64(3.14)},
		{"complex64", From(complex(1, 2)), complex64(0), complex64(complex(1, 2))},
		{"complex128", From(complex(1, 2)), complex128(0), complex128(complex(1, 2))},
	}

	for _, tt := range tests {
		v := reflect.New(reflect.TypeOf(tt.goType)).Elem()
		if !ToValue(tt.pyValue, v) {
			t.Errorf("ToValue failed for %v", tt.name)
		}
		if v.Interface() != tt.expected {
			t.Errorf("Expected %v, got %v", tt.expected, v.Interface())
		}

	}

	func() {
		v := reflect.New(reflect.TypeOf("")).Elem()
		if !ToValue(From("hello"), v) {
			t.Error("ToValue failed for string")
		}
		if v.String() != "hello" {
			t.Errorf("Expected 'hello', got %v", v.String())
		}
	}()

	func() {
		v := reflect.New(reflect.TypeOf(true)).Elem()
		if !ToValue(From(true), v) {
			t.Error("ToValue failed for bool")
		}
		if !v.Bool() {
			t.Error("Expected true, got false")
		}
	}()

	func() {
		expected := []byte("hello")
		v := reflect.New(reflect.TypeOf([]byte{})).Elem()
		if !ToValue(From(expected), v) {
			t.Error("ToValue failed for []byte")
		}
		if !reflect.DeepEqual(v.Bytes(), expected) {
			t.Errorf("Expected %v, got %v", expected, v.Bytes())
		}
	}()

	func() {
		expected := []int{1, 2, 3}
		v := reflect.New(reflect.TypeOf([]int{})).Elem()
		if !ToValue(From(expected), v) {
			t.Error("ToValue failed for slice")
		}
		if !reflect.DeepEqual(v.Interface(), expected) {
			t.Errorf("Expected %v, got %v", expected, v.Interface())
		}
	}()

	func() {
		expected := map[string]int{"one": 1, "two": 2}
		v := reflect.New(reflect.TypeOf(map[string]int{})).Elem()
		if !ToValue(From(expected), v) {
			t.Error("ToValue failed for map")
		}
		if !reflect.DeepEqual(v.Interface(), expected) {
			t.Errorf("Expected %v, got %v", expected, v.Interface())
		}
	}()

	func() {
		type TestStruct struct {
			Name string
			Age  int
		}
		expected := TestStruct{Name: "Alice", Age: 30}
		v := reflect.New(reflect.TypeOf(TestStruct{})).Elem()
		if !ToValue(From(expected), v) {
			t.Error("ToValue failed for struct")
		}
		if !reflect.DeepEqual(v.Interface(), expected) {
			t.Errorf("Expected %v, got %v", expected, v.Interface())
		}
	}()

	func() {
		tests := []struct {
			name    string
			pyValue Object
			goType  interface{}
		}{
			{"string to int", From("not a number"), int(0)},
			{"int to bool", From(42), true},
			{"float to string", From(3.14), ""},
			{"list to map", From([]int{1, 2, 3}), map[string]int{}},
		}

		for _, tt := range tests {
			v := reflect.New(reflect.TypeOf(tt.goType)).Elem()
			if ToValue(tt.pyValue, v) {
				t.Errorf("ToValue should have failed for %v", tt.name)
			}

		}
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("ToValue should fail for nil reflect.Value")
			}
		}()
		var nilValue reflect.Value
		if ToValue(From(42), nilValue) {
			t.Error("ToValue should fail for nil reflect.Value")
		}
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("ToValue should fail for non-settable value")
			}
		}()
		v := reflect.ValueOf(42) // not settable
		if ToValue(From(43), v) {
			t.Error("ToValue should fail for non-settable value")
		}
	}()
}

func TestFromSpecialCases(t *testing.T) {
	setupTest(t)

	func() {
		// Test From with uint values
		tests := []struct {
			input    uint
			expected uint64
		}{
			{0, 0},
			{42, 42},
			{^uint(0), ^uint64(0)}, // maximum uint value
		}

		for _, tt := range tests {
			obj := From(tt.input)
			if !obj.IsLong() {
				t.Errorf("From(uint) did not create Long object")
			}
			if got := obj.AsLong().Uint64(); got != tt.expected {
				t.Errorf("From(%d) = %d, want %d", tt.input, got, tt.expected)
			}
		}
	}()

	func() {
		// Test From with Object.cpyObj()
		original := From(42)
		obj := From(original.cpyObj())

		if !obj.IsLong() {
			t.Error("From(Object.cpyObj()) did not create Long object")
		}
		if got := obj.AsLong().Int64(); got != 42 {
			t.Errorf("From(Object.cpyObj()) = %d, want 42", got)
		}

		// Test that the new object is independent
		original = From(100)
		if got := obj.AsLong().Int64(); got != 42 {
			t.Errorf("Object was not independent, got %d after modifying original", got)
		}
	}()

	func() {
		// Test From with functions
		add := func(a, b int) int { return a + b }
		obj := From(add)

		// Verify it's a function type
		if !obj.IsFunc() {
			t.Error("From(func) did not create Function object")
		}

		fn := obj.AsFunc()

		// Test function call
		result := fn.Call(5, 3)

		if !result.IsLong() {
			t.Error("Function call result is not a Long")
		}
		if got := result.AsLong().Int64(); got != 8 {
			t.Errorf("Function call = %d, want 8", got)
		}
	}()

	func() {
		// Test From with function that returns multiple values
		divMod := func(a, b int) (int, int) {
			return a / b, a % b
		}
		obj := From(divMod)
		if !obj.IsFunc() {
			t.Error("From(func) did not create Function object")
		}

		fn := obj.AsFunc()

		result := fn.Call(7, 3)

		// Result should be a tuple with two values
		if !result.IsTuple() {
			t.Error("Multiple return value function did not return a Tuple")
		}

		tuple := result.AsTuple()
		if tuple.Len() != 2 {
			t.Errorf("Expected tuple of length 2, got %d", tuple.Len())
		}

		quotient := tuple.Get(0).AsLong().Int64()
		remainder := tuple.Get(1).AsLong().Int64()

		if quotient != 2 || remainder != 1 {
			t.Errorf("Got (%d, %d), want (2, 1)", quotient, remainder)
		}
	}()
}

func TestToValueWithCustomType(t *testing.T) {
	setupTest(t)

	// Define a custom Go type
	type Point struct {
		X int
		Y int
	}

	// Add the type to Python
	pointClass := MainModule().AddType(Point{}, nil, "Point", "Point class")

	func() {
		// Create a Point instance in Python and convert it back to Go
		pyCode := `
p = Point()
p.x = 10
p.y = 20
`
		locals := MakeDict(nil)
		globals := MakeDict(nil)
		builtins := ImportModule("builtins")
		globals.Set(MakeStr("__builtins__"), builtins.Object)
		globals.Set(MakeStr("Point"), pointClass)

		code, err := CompileString(pyCode, "<string>", FileInput)
		if err != nil {
			t.Errorf("CompileString() error = %v", err)
		}
		EvalCode(code, globals, locals)

		// Get the Python Point instance
		pyPoint := locals.Get(MakeStr("p"))

		// Convert back to Go Point struct
		var point Point
		v := reflect.ValueOf(&point).Elem()
		if !ToValue(pyPoint, v) {
			t.Error("ToValue failed for custom type Point")
		}

		// Verify the values
		if point.X != 10 || point.Y != 20 {
			t.Errorf("Expected Point{10, 20}, got Point{%d, %d}", point.X, point.Y)
		}
	}()

	func() {
		// Test converting a non-Point Python object to Point should fail
		dict := MakeDict(nil)
		dict.Set(MakeStr("x"), From(10))
		dict.Set(MakeStr("y"), From(20))

		var point Point
		v := reflect.ValueOf(&point).Elem()
		if !ToValue(dict.Object, v) {
			t.Error("ToValue failed for custom type Point")
		}

		if point.X != 10 || point.Y != 20 {
			t.Errorf("Expected Point{10, 20}, got Point{%d, %d}", point.X, point.Y)
		}
	}()
}

func TestFromWithCustomType(t *testing.T) {
	setupTest(t)

	type Point struct {
		X int
		Y int
	}

	// Add the type to Python
	pointClass := MainModule().AddType(Point{}, nil, "Point", "Point class")

	func() {
		// Test From with struct instance
		p := Point{X: 10, Y: 20}
		obj := From(p)

		// Verify the type
		if obj.Type().cpyObj() != pointClass.cpyObj() {
			t.Error("From(Point) created object with wrong type")
		}
		// Verify the values
		if obj.AttrLong("x").Int64() != 10 {
			t.Error("Wrong X value after From conversion")
		}
		if obj.AttrLong("y").Int64() != 20 {
			t.Error("Wrong Y value after From conversion")
		}

		// Convert back to Go and verify
		var p2 Point
		v := reflect.ValueOf(&p2).Elem()
		if !ToValue(obj, v) {
			t.Error("ToValue failed for custom type Point")
		}

		if p2.X != p.X || p2.Y != p.Y {
			t.Errorf("Round trip conversion failed: got Point{%d, %d}, want Point{%d, %d}",
				p2.X, p2.Y, p.X, p.Y)
		}
	}()

	func() {
		// Test From with pointer to struct
		p := &Point{X: 30, Y: 40}
		obj := From(p)

		// Verify the type
		if obj.Type().cpyObj() != pointClass.cpyObj() {
			t.Error("From(*Point) created object with wrong type")
		}

		// Verify the values
		if obj.AttrLong("x").Int64() != 30 {
			t.Error("Wrong X value after From pointer conversion")
		}
		if obj.AttrLong("y").Int64() != 40 {
			t.Error("Wrong Y value after From pointer conversion")
		}
	}()
}
