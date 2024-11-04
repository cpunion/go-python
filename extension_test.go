package gp

import (
	"testing"
)

// TestStruct contains various types of fields for testing
type TestStruct struct {
	// C-compatible basic types
	BoolField       bool
	Int8Field       int8
	Int16Field      int16
	Int32Field      int32
	Int64Field      int64
	IntField        int
	Uint8Field      uint8
	Uint16Field     uint16
	Uint32Field     uint32
	Uint64Field     uint64
	UintField       uint
	Float32Field    float32
	Float64Field    float64
	Complex64Field  complex64
	Complex128Field complex128

	// Non-C-compatible types
	StringField string
	SliceField  []int
	MapField    map[string]int
	StructField struct{ X int }
}

func (t *TestStruct) TestMethod() int {
	return 42
}

func TestAddType(t *testing.T) {
	setupTest(t)
	m := MainModule()

	// test add type
	typ := m.AddType(TestStruct{}, nil, "TestStruct", "Test struct documentation")
	if typ.Nil() {
		t.Fatal("Failed to create type")
	}

	// test type by Python code
	code := `
# create instance
obj = TestStruct()

# test C-compatible types
obj.bool_field = True
obj.int8_field = 127
obj.int16_field = 32767
obj.int32_field = 2147483647
obj.int64_field = 9223372036854775807
obj.int_field = 1234567890
obj.uint8_field = 255
obj.uint16_field = 65535
obj.uint32_field = 4294967295
obj.uint64_field = 18446744073709551615
obj.uint_field = 4294967295
obj.float32_field = 3.14
obj.float64_field = 3.14159265359
obj.complex64_field = 1.5 + 2.5j
obj.complex128_field = 3.14 + 2.718j

# test non-C-compatible types
obj.string_field = "test string"
assert obj.string_field == "test string"

obj.slice_field = [1, 2, 3]
assert obj.slice_field == [1, 2, 3]

obj.map_field = {"key": 42}
assert obj.map_field["key"] == 42

obj.struct_field = {"x": 100}
assert obj.struct_field.x == 100

# test method call
result = obj.test_method()
assert result == 42

# verify C-compatible types
assert obj.bool_field == True
assert obj.int8_field == 127
assert obj.int16_field == 32767
assert obj.int32_field == 2147483647
assert obj.int64_field == 9223372036854775807
assert obj.int_field == 1234567890
assert obj.uint8_field == 255
assert obj.uint16_field == 65535
assert obj.uint32_field == 4294967295
assert obj.uint64_field == 18446744073709551615
assert obj.uint_field == 4294967295
assert abs(obj.float32_field - 3.14) < 0.0001
assert abs(obj.float64_field - 3.14159265359) < 0.0000001
assert abs(obj.complex64_field - (1.5 + 2.5j)) < 0.0001
assert abs(obj.complex128_field - (3.14 + 2.718j)) < 0.0000001

# verify non-C-compatible types
assert obj.string_field == "test string"
assert obj.slice_field == [1, 2, 3]
assert obj.map_field["key"] == 42
assert obj.struct_field.x == 100
`

	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

type InitTestStruct struct {
	Value int
}

func (i *InitTestStruct) Init(val int) {
	i.Value = val
}

func AddTypeWithInit(t *testing.T) {
	setupTest(t)
	m := MainModule()

	typ := m.AddType(InitTestStruct{}, (*InitTestStruct).Init, "InitTestStruct", "Test init struct")
	if typ.Nil() {
		t.Fatal("Failed to create type with init")
	}

	// test init function
	code := `
# test init function
obj = InitTestStruct(42)
assert obj.value == 42

# test error handling without arguments
try:
    obj2 = InitTestStruct()
    assert False, "Should fail without arguments"
except TypeError as e:
    pass

# test error handling with wrong argument type
try:
    obj3 = InitTestStruct("wrong type")
    assert False, "Should fail with wrong argument type"
except TypeError:
    pass
`

	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestCreateFunc(t *testing.T) {
	setupTest(t)

	// Test simple function
	simpleFunc := func(x int) int {
		return x * 2
	}
	f1 := CreateFunc("simple_func", simpleFunc, "Doubles the input value")
	if f1.Nil() {
		t.Fatal("Failed to create simple function")
	}

	// Test function with multiple arguments and return values
	multiFunc := func(x int, y string) (int, string) {
		return x * 2, y + y
	}
	f2 := CreateFunc("multi_func", multiFunc, "Returns doubled number and duplicated string")
	if f2.Nil() {
		t.Fatal("Failed to create function with multiple returns")
	}

	// Test the functions using Python code
	code := `
# Test simple function
result = simple_func(21)
assert result == 42, f"Expected 42, got {result}"

# Test multiple return values
num, text = multi_func(5, "hello")
assert num == 10, f"Expected 10, got {num}"
assert text == "hellohello", f"Expected 'hellohello', got {text}"

# Test error handling - wrong argument type
try:
    simple_func("not a number")
    assert False, "Should fail with wrong argument type"
except TypeError:
    pass

# Test error handling - wrong number of arguments
try:
    simple_func(1, 2)
    assert False, "Should fail with wrong number of arguments"
except TypeError:
    pass
`

	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	// Test invalid function type
	defer func() {
		if r := recover(); r == nil {
			t.Error("CreateFunc should panic with non-function argument")
		}
	}()
	CreateFunc("non_func", 42, "This should panic")
}

func explicitFunc(x int) int {
	return x + 1
}

func namedFunc(x string) string {
	return "Hello " + x
}

func TestModuleAddMethod(t *testing.T) {
	setupTest(t)
	m := MainModule()

	// Test with explicit name
	f1 := m.AddMethod("", explicitFunc, " - adds one to input")
	if f1.Nil() {
		t.Fatal("Failed to create function with explicit name")
	}

	// Test with empty name (should use function name)
	f2 := m.AddMethod("", namedFunc, " - adds greeting")
	if f2.Nil() {
		t.Fatal("Failed to create function with derived name")
	}

	// Test with anonymous function (should generate name)
	f3 := m.AddMethod("", func(x, y int) int {
		return x * y
	}, " - multiplies two numbers")
	if f3.Nil() {
		t.Fatal("Failed to create anonymous function")
	}

	code := `
# Test explicit named function
result = explicit_func(41)
assert result == 42, f"Expected 42, got {result}"

# Test function with derived name
result = named_func("World")
assert result == "Hello World", f"Expected 'Hello World', got {result}"

# Test documentation
import sys
if sys.version_info >= (3, 2):
    assert explicit_func.__doc__.strip() == "explicit_func - adds one to input"
    assert named_func.__doc__.strip() == "named_func - adds greeting"

# Test error cases
try:
    explicit_func("wrong type")
    assert False, "Should fail with wrong argument type"
except TypeError:
    pass

try:
    explicit_func(1, 2)
    assert False, "Should fail with wrong number of arguments"
except TypeError:
    pass
`

	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	// Test invalid function type
	defer func() {
		if r := recover(); r == nil {
			t.Error("AddMethod should panic with non-function argument")
		}
	}()
	m.AddMethod("invalid", "not a function", "")
}

func TestAddTypeWithPtrReceiverInit(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type InitTestType struct {
		Value  int
		Name   string
		Active bool
	}

	ptrInit := func(t *InitTestType, val int, name string, active bool) {
		t.Value = val
		t.Name = name
		t.Active = active
	}
	typ := m.AddType(InitTestType{}, ptrInit, "InitTestType", "")
	if typ.Nil() {
		t.Fatal("Failed to create type with pointer receiver init")
	}

	code := `
# Test pointer receiver init with multiple args
obj = InitTestType(42, "hello", True)
assert obj.value == 42
assert obj.name == "hello"
assert obj.active == True

# Test error handling
try:
    obj2 = InitTestType(42)  # Missing arguments
    assert False, "Should fail with wrong number of arguments"
except TypeError:
    pass

try:
    obj3 = InitTestType("wrong", "type", True)  # Wrong argument type
    assert False, "Should fail with wrong argument type"
except TypeError:
    pass
`
	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestAddTypeWithValueConstructor(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type InitTestType struct {
		Value  int
		Name   string
		Active bool
	}

	constructorInit := func(val int, name string, active bool) InitTestType {
		return InitTestType{
			Value:  val,
			Name:   name,
			Active: active,
		}
	}
	typ := m.AddType(InitTestType{}, constructorInit, "InitTestType", "")
	if typ.Nil() {
		t.Fatal("Failed to create type with value constructor")
	}

	code := `
# Test value constructor with multiple args
obj = InitTestType(43, "world", False)
assert obj.value == 43
assert obj.name == "world"
assert obj.active == False
`
	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func TestAddTypeWithPtrConstructor(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type InitTestType struct {
		Value  int
		Name   string
		Active bool
	}

	ptrConstructorInit := func(val int, name string) *InitTestType {
		return &InitTestType{
			Value: val,
			Name:  name,
		}
	}
	typ := m.AddType(InitTestType{}, ptrConstructorInit, "InitTestType", "")
	if typ.Nil() {
		t.Fatal("Failed to create type with pointer constructor")
	}

	code := `
# Test pointer constructor with multiple args
obj = InitTestType(44, "python")
assert obj.value == 44
assert obj.name == "python"
assert obj.active == False  # default value
`
	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

type Inner struct {
	X int
	Y string
}

func TestAddTypeRecursive(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type Outer struct {
		Inner     Inner
		InnerPtr  *Inner
		Value     int
		InnerList []Inner
	}

	// Register Outer type - should automatically register Inner
	obj := m.AddType(Outer{}, nil, "Outer", "")
	if obj.Nil() {
		t.Fatal("Failed to create Outer type")
	}

	code := `
# Test nested struct access
o = Outer()
o.inner.x = 42
o.inner.y = "hello"
assert o.inner.x == 42
assert o.inner.y == "hello"

# Test pointer to struct
o.inner_ptr = Inner()
o.inner_ptr.x = 43
o.inner_ptr.y = "world"
assert o.inner_ptr.x == 43
assert o.inner_ptr.y == "world"

# Test basic field
o.value = 100
assert o.value == 100

# Test slice of structs
o.inner_list = [Inner()]
o.inner_list[0].x = 44
o.inner_list[0].y = "python"
assert o.inner_list[0].x == 44
assert o.inner_list[0].y == "python"
`

	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}