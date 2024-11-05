package gp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestObjectCreation(t *testing.T) {
	setupTest(t)
	// Test From() with different Go types
	tests := []struct {
		name     string
		input    interface{}
		checkFn  func(Object) bool
		expected interface{}
	}{
		{"int", 42, func(o Object) bool { return o.IsLong() }, 42},
		{"float64", 3.14, func(o Object) bool { return o.IsFloat() }, 3.14},
		{"string", "hello", func(o Object) bool { return o.IsStr() }, "hello"},
		{"bool", true, func(o Object) bool { return o.IsBool() }, true},
		{"[]byte", []byte("bytes"), func(o Object) bool { return o.IsBytes() }, []byte("bytes")},
		{"slice", []int{1, 2, 3}, func(o Object) bool { return o.IsList() }, []int{1, 2, 3}},
		{"map", map[string]int{"a": 1}, func(o Object) bool { return o.IsDict() }, map[string]int{"a": 1}},
	}

	for _, tt := range tests {
		obj := From(tt.input)
		if !tt.checkFn(obj) {
			t.Errorf("From(%v) created wrong type", tt.input)
		}

		// Test conversion back to Go value
		switch expected := tt.expected.(type) {
		case int:
			if got := obj.AsLong().Int64(); got != int64(expected) {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		case float64:
			if got := obj.AsFloat().Float64(); got != expected {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		case string:
			if got := obj.AsStr().String(); got != expected {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		case bool:
			if got := obj.AsBool().Bool(); got != expected {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		case []byte:
			if got := obj.AsBytes().Bytes(); !reflect.DeepEqual(got, expected) {
				t.Errorf("Expected %v, got %v", expected, got)
			}
		}
	}
}

func TestObjectAttributes(t *testing.T) {
	setupTest(t)
	// Test attributes using Python's built-in object type
	builtins := ImportModule("builtins")
	obj := builtins.AttrFunc("object").Call()

	// Get built-in attribute
	cls := obj.Attr("__class__")
	if cls.Nil() {
		t.Error("Failed to get __class__ attribute")
	}

	// Test Dir() method
	dir := obj.Dir()
	if dir.Len() == 0 {
		t.Error("Dir() returned empty list")
	}

	// Create a custom class to test attribute setting
	pyCode := `
class TestClass:
    pass
`
	locals := MakeDict(nil)
	globals := MakeDict(nil)
	globals.Set(MakeStr("__builtins__"), builtins.Object)

	code, err := CompileString(pyCode, "<string>", FileInput)
	if err != nil {
		t.Errorf("CompileString() error = %v", err)
	}

	EvalCode(code, globals, locals).AsModule()
	testClass := locals.Get(MakeStr("TestClass")).AsFunc()
	instance := testClass.Call()

	// Now we can set attributes
	instance.SetAttr("new_attr", "test_value")
	value := instance.Attr("new_attr")
	if value.AsStr().String() != "test_value" {
		t.Error("SetAttr failed to set new attribute")
	}
}

func TestDictOperations(t *testing.T) {
	setupTest(t)
	// Test dictionary operations
	pyDict := MakeDict(nil)
	pyDict.Set(MakeStr("key1"), From(42))
	pyDict.Set(MakeStr("key2"), From("value"))

	value := pyDict.Get(MakeStr("key1"))
	if value.AsLong().Int64() != 42 {
		t.Error("Failed to get dictionary item")
	}

	func() {
		pyDict.Set(MakeStr("key3"), From("new_value"))
		value := pyDict.Get(MakeStr("key3"))
		if value.AsStr().String() != "new_value" {
			t.Error("Failed to set dictionary item")
		}
	}()
}

func TestObjectConversion(t *testing.T) {
	setupTest(t)
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "Alice", Age: 30}
	obj := From(person)

	if !obj.IsDict() {
		t.Error("Struct should be converted to Python dict")
	}

	dict := obj.AsDict()
	if dict.Get(From("name")).AsStr().String() != "Alice" {
		t.Error("Failed to convert struct field 'Name'")
	}
	if dict.Get(From("age")).AsLong().Int64() != 30 {
		t.Error("Failed to convert struct field 'Age'")
	}

	func() {
		slice := []int{1, 2, 3}
		obj := From(slice)

		if !obj.IsList() {
			t.Error("Slice should be converted to Python list")
		}

		list := obj.AsList()
		if list.Len() != 3 {
			t.Error("Wrong length after conversion")
		}
		if list.GetItem(0).AsLong().Int64() != 1 {
			t.Error("Wrong value at index 0")
		}
	}()
}

func TestObjectString(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"int", 42, "42"},
		{"string", "hello", "hello"},
		{"bool", true, "True"},
	}

	for _, tt := range tests {
		obj := From(tt.input)
		str := obj.String()
		if str != tt.expected {
			t.Errorf("String() = %v, want %v", str, tt.expected)
		}

	}
}

func TestPyObjectMethods(t *testing.T) {
	setupTest(t)
	// Test pyObject.Obj()
	obj := From(42)
	if obj.pyObject.Obj() == nil {
		t.Error("pyObject.Obj() returned nil for valid object")
	}

	func() {
		var nilObj *pyObject
		if nilObj.Obj() != nil {
			t.Error("pyObject.Obj() should return nil for nil object")
		}
	}()

	func() {
		// Test pyObject.Ensure()
		obj := From(42)
		obj.Ensure() // Should not panic
	}()

	func() {
		var nilObj Object
		defer func() {
			if r := recover(); r == nil {
				t.Error("Ensure() should panic for nil object")
			}
		}()
		nilObj.Ensure()
	}()
}

func TestObjectMethods(t *testing.T) {
	setupTest(t)
	// Test Object.object()
	obj := From(42)
	if obj.object() != obj {
		t.Error("object() should return the same object")
	}

	// Test Object.Attr* methods
	// Create a test class with various attribute types
	pyCode := `
class TestClass:
    def __init__(self):
        self.int_val = 42
        self.float_val = 3.14
        self.str_val = "test"
        self.bool_val = True
        self.list_val = [1, 2, 3]
        self.dict_val = {"key": "value"}
        self.tuple_val = (1, 2, 3)
`
	locals := MakeDict(nil)
	globals := MakeDict(nil)
	builtins := ImportModule("builtins")
	globals.Set(MakeStr("__builtins__"), builtins.Object)

	code, err := CompileString(pyCode, "<string>", FileInput)
	if err != nil {
		t.Errorf("CompileString() error = %v", err)
	}
	EvalCode(code, globals, locals)

	testClass := locals.Get(MakeStr("TestClass")).AsFunc()
	instance := testClass.Call()

	// Test each Attr* method
	if instance.AttrLong("int_val").Int64() != 42 {
		t.Error("AttrLong failed")
	}
	if instance.AttrFloat("float_val").Float64() != 3.14 {
		t.Error("AttrFloat failed")
	}
	if instance.AttrString("str_val").String() != "test" {
		t.Error("AttrString failed")
	}
	if !instance.AttrBool("bool_val").Bool() {
		t.Error("AttrBool failed")
	}
	if instance.AttrList("list_val").Len() != 3 {
		t.Error("AttrList failed")
	}
	if instance.AttrDict("dict_val").Get(MakeStr("key")).AsStr().String() != "value" {
		t.Error("AttrDict failed")
	}
	if instance.AttrTuple("tuple_val").Len() != 3 {
		t.Error("AttrTuple failed")
	}

	func() {
		// Test Object.IsTuple and AsTuple
		// Create a Python tuple using Python code to ensure proper tuple creation
		pyCode := `
def make_tuple():
    return (1, 2, 3)
`
		locals := MakeDict(nil)
		globals := MakeDict(nil)
		builtins := ImportModule("builtins")
		globals.Set(MakeStr("__builtins__"), builtins.Object)

		code, err := CompileString(pyCode, "<string>", FileInput)
		if err != nil {
			t.Errorf("CompileString() error = %v", err)
		}
		EvalCode(code, globals, locals)

		makeTuple := locals.Get(MakeStr("make_tuple")).AsFunc()
		tuple := makeTuple.Call()

		// Test IsTuple
		if !tuple.IsTuple() {
			t.Error("IsTuple failed to identify tuple")
		}

		// Test AsTuple
		pythonTuple := tuple.AsTuple()
		if pythonTuple.Len() != 3 {
			t.Error("AsTuple conversion failed")
		}

		// Verify tuple contents
		if pythonTuple.Get(0).AsLong().Int64() != 1 {
			t.Error("Incorrect value at index 0")
		}
		if pythonTuple.Get(1).AsLong().Int64() != 2 {
			t.Error("Incorrect value at index 1")
		}
		if pythonTuple.Get(2).AsLong().Int64() != 3 {
			t.Error("Incorrect value at index 2")
		}
	}()

	func() {
		// Test Object.Repr and Type
		obj := From(42)
		if obj.Repr() != "42" {
			t.Error("Repr failed")
		}
	}()

	func() {
		typeObj := obj.Type()
		if typeObj.Repr() != "<class 'int'>" {
			t.Error("Type failed")
		}
	}()

	func() {
		// Test From with various numeric types
		tests := []struct {
			input    interface{}
			expected int64
		}{
			{int8(42), 42},
			{int16(42), 42},
			{int32(42), 42},
			{int64(42), 42},
			{uint8(42), 42},
			{uint16(42), 42},
			{uint32(42), 42},
			{uint64(42), 42},
		}

		for _, tt := range tests {
			obj := From(tt.input)
			if obj.AsLong().Int64() != tt.expected {
				t.Errorf("From(%T) = %v, want %v", tt.input, obj.AsLong().Int64(), tt.expected)
			}
		}
	}()

	func() {
		// Test From with false boolean
		obj := From(false)
		if obj.AsBool().Bool() != false {
			t.Error("From(false) failed")
		}
	}()

	func() {
		// Test Object.Obj()
		obj := From(42)
		if obj.Obj() == nil {
			t.Error("Object.Obj() returned nil for valid object")
		}
	}()

	func() {
		var nilObj Object
		if nilObj.Obj() != nil {
			t.Error("Object.Obj() should return nil for nil object")
		}
	}()

	func() {
		// Test AttrBytes
		builtins := ImportModule("types")
		objType := builtins.AttrFunc("SimpleNamespace")
		obj := objType.Call()

		// Create a simple object with bytes attribute
		obj.SetAttr("bytes_val", From([]byte("hello")))

		if !bytes.Equal(obj.AttrBytes("bytes_val").Bytes(), []byte("hello")) {
			t.Error("AttrBytes failed")
		}
	}()

	func() {
		// Test Object.Call with kwargs
		pyCode := `
def test_func(a, b=10, c="default"):
    return (a, b, c)
`
		locals := MakeDict(nil)
		globals := MakeDict(nil)
		globals.Set(MakeStr("__builtins__"), builtins.Object)

		code, err := CompileString(pyCode, "<string>", FileInput)
		if err != nil {
			t.Errorf("CompileString() error = %v", err)
		}
		EvalCode(code, globals, locals)

		testFunc := locals.Get(MakeStr("test_func"))

		// Call with positional and keyword arguments
		result := testFunc.Call("__call__", 1, KwArgs{
			"b": 20,
			"c": "custom",
		})

		tuple := result.AsTuple()
		if tuple.Get(0).AsLong().Int64() != 1 {
			t.Error("Wrong value for first argument")
		}
		if tuple.Get(1).AsLong().Int64() != 20 {
			t.Error("Wrong value for keyword argument b")
		}
		if tuple.Get(2).AsStr().String() != "custom" {
			t.Error("Wrong value for keyword argument c")
		}
	}()
}

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
		// Test From with Object.Obj()
		original := From(42)
		obj := From(original.Obj())

		if !obj.IsLong() {
			t.Error("From(Object.Obj()) did not create Long object")
		}
		if got := obj.AsLong().Int64(); got != 42 {
			t.Errorf("From(Object.Obj()) = %d, want 42", got)
		}

		// Test that the new object is independent
		original = From(100)
		if got := obj.AsLong().Int64(); got != 42 {
			t.Errorf("Object was not independent, got %d after modifying original", got)
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
		if obj.Type().Obj() != pointClass.Obj() {
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
		if obj.Type().Obj() != pointClass.Obj() {
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
