package gp

import (
	"fmt"
	"strings"
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
assert obj.uint16_field == 65535, f"Expected 65535, got {obj.uint16_field}"
assert obj.uint32_field == 4294967295, f"Expected 4294967295, got {obj.uint32_field}"
assert obj.uint64_field == 18446744073709551615, f"Expected 18446744073709551615, got {obj.uint64_field}"
assert obj.uint_field == 4294967295, f"Expected 4294967295, got {obj.uint_field}"
assert abs(obj.float32_field - 3.14) < 0.0001, f"Expected 3.14, got {obj.float32_field}"
assert abs(obj.float64_field - 3.14159265359) < 0.0000001, f"Expected 3.14159265359, got {obj.float64_field}"
assert abs(obj.complex64_field - (1.5 + 2.5j)) < 0.0001, f"Expected (1.5 + 2.5j), got {obj.complex64_field}"
assert abs(obj.complex128_field - (3.14 + 2.718j)) < 0.0000001, f"Expected (3.14 + 2.718j), got {obj.complex128_field}"

# verify non-C-compatible types
assert obj.string_field == "test string", f"Expected 'test string', got {obj.string_field}"
assert obj.slice_field == [1, 2, 3], f"Expected [1, 2, 3], got {obj.slice_field}"
assert obj.map_field["key"] == 42, f"Expected 42, got {obj.map_field['key']}"
assert obj.struct_field.x == 100, f"Expected 100, got {obj.struct_field.x}"
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

func TestAddTypeWithInit(t *testing.T) {
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

	tests := []struct {
		name     string
		fn       any
		doc      string
		testCode string
	}{
		{
			name: "simple",
			fn: func(x int) int {
				return x * 2
			},
			doc: "Doubles the input value",
			testCode: `
result = simple_func(21)
assert result == 42, f"Expected 42, got {result}"
assert str(inspect.signature(simple_func)) == "(arg0, /)"
`,
		},
		{
			name: "multi_args",
			fn: func(x int, y string) (int, string) {
				return x * 2, y + y
			},
			doc: "Returns doubled number and duplicated string",
			testCode: `
num, text = multi_args_func(5, "hello")
assert num == 10, f"Expected 10, got {num}"
assert text == "hellohello", f"Expected 'hellohello', got {text}"
assert str(inspect.signature(multi_args_func)) == "(arg0, arg1, /)"
`,
		},
		{
			name: "with_kwargs",
			fn: func(name string, kwargs KwArgs) Object {
				return None()
			},
			doc: "Function with kwargs",
			testCode: `
result = with_kwargs_func("test", extra="value")
assert result is None
assert str(inspect.signature(with_kwargs_func)) == "(arg0, /, **kwargs)"
`,
		},
		{
			name: "no_args",
			fn: func() string {
				return "hello"
			},
			doc: "Function with no arguments",
			testCode: `
result = no_args_func()
assert result == "hello"
assert str(inspect.signature(no_args_func)) == "()"
`,
		},
		{
			name: "only_kwargs",
			fn: func(kwargs KwArgs) Object {
				return None()
			},
			doc: "Function with only kwargs",
			testCode: `
result = only_kwargs_func(x=1, y=2)
assert result is None
assert str(inspect.signature(only_kwargs_func)) == "(**kwargs)"
`,
		},
	}

	code := `
import inspect

`
	for _, tt := range tests {
		funcName := tt.name + "_func"
		f := CreateFunc(funcName, tt.fn, tt.doc)
		if f.Nil() {
			t.Fatalf("Failed to create function for test case: %s", tt.name)
		}
		code += tt.testCode + "\n"
	}
	t.Log(code)
	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	// Test failure cases
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("CreateFunc should panic with non-function argument")
			}
		}()
		CreateFunc("invalid", 42, "should panic")
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("CreateFunc should panic with nil function")
			}
		}()
		CreateFunc("nil_func", nil, "should panic")
	}()
}

func TestCreateFuncInvalid(t *testing.T) {
	setupTest(t)
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

func TestModuleAddMethod(t *testing.T) {
	setupTest(t)
	m := MainModule()

	tests := []struct {
		name     string
		fn       any
		doc      string
		testCode string
	}{
		{
			name: "explicit",
			fn:   explicitFunc,
			doc:  "adds one to input",
			testCode: `
result = explicit_func(41)
assert result == 42, f"Expected 42, got {result}"
assert str(inspect.signature(explicit_func)) == "(arg0, /)"
`,
		},
		{
			name: "with_kwargs",
			fn: func(x int, kwargs KwArgs) int {
				return x
			},
			doc: "function with kwargs",
			testCode: `
result = with_kwargs_func(42, extra="value")
assert result == 42
assert str(inspect.signature(with_kwargs_func)) == "(arg0, /, **kwargs)"
`,
		},
		{
			name: "multi_args",
			fn: func(x, y int) int {
				return x * y
			},
			doc: "multiplies two numbers",
			testCode: `
result = multi_args_func(6, 7)
assert result == 42
assert str(inspect.signature(multi_args_func)) == "(arg0, arg1, /)"
`,
		},
		{
			name: "no_args",
			fn: func() string {
				return "hello"
			},
			doc: "returns hello",
			testCode: `
result = no_args_func()
assert result == "hello"
assert str(inspect.signature(no_args_func)) == "()"
`,
		},
		{
			name: "only_kwargs",
			fn: func(kwargs KwArgs) Object {
				return None()
			},
			doc: "function with only kwargs",
			testCode: `
result = only_kwargs_func(x=1, y=2)
assert result is None
assert str(inspect.signature(only_kwargs_func)) == "(**kwargs)"
`,
		},
	}

	code := `
import inspect

`
	for _, tt := range tests {
		funcName := tt.name + "_func"
		f := m.AddMethod(funcName, tt.fn, tt.doc)
		if f.Nil() {
			t.Fatalf("Failed to create function for test case: %s", tt.name)
		}
		code += tt.testCode + "\n"
	}

	t.Log(code)
	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("AddMethod should panic with non-function argument")
			}
		}()
		m.AddMethod("invalid", 42, "should panic")
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("AddMethod should panic with nil function")
			}
		}()
		m.AddMethod("nil_func", nil, "should panic")
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("AddMethod should panic with empty module")
			}
		}()
		var emptyModule Module
		emptyModule.AddMethod("empty", func() {}, "should panic")
	}()
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

func TestSetterMethodEdgeCases(t *testing.T) {
	setupTest(t)

	type ChildStruct struct {
		Value int
	}

	type ParentStruct struct {
		unexported int
		Value      int
		Child      *ChildStruct
		Nested     ChildStruct
	}

	m := MainModule()
	m.AddType(ChildStruct{}, nil, "ChildStruct", "")
	m.AddType(ParentStruct{}, nil, "ParentStruct", "")

	code := `
obj = ParentStruct()
try:
    obj.value = "invalid"  # Try to set int with string
    assert False, "Should have raised TypeError"
except TypeError:
    pass

try:
    obj.child = 123  # Try to set struct pointer with int
    assert False, "Should have raised TypeError"
except TypeError:
    pass

try:
    obj.nested = 123  # Try to set struct with int
    assert False, "Should have raised TypeError"
except TypeError:
    pass
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetterMethodEdgeCases(t *testing.T) {
	setupTest(t)

	type ChildStruct struct {
		Value int
	}

	type ParentStruct struct {
		Value  int
		Child  *ChildStruct
		Nested ChildStruct
	}

	m := MainModule()
	m.AddType(ChildStruct{}, nil, "ChildStruct", "")
	m.AddType(ParentStruct{}, nil, "ParentStruct", "")

	code := `
obj = ParentStruct()
obj.child = None  # Set pointer to nil
val = obj.child   # Should return None for nil pointer
assert val is None

obj.nested = ChildStruct()  # Set nested struct
val = obj.nested  # Should return wrapper for nested struct
assert isinstance(val, ChildStruct)

# Test accessing nested struct fields
obj.nested.value = 42
assert obj.nested.value == 42
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWrapperMethodEdgeCases(t *testing.T) {
	setupTest(t)

	type TestStruct struct {
		Value int
	}

	m := MainModule()

	// Test method with wrong number of arguments
	m.AddMethod("test_func", func(x int, y int) int { return x + y }, "")

	code := `
try:
    test_func(1)  # Missing argument
    assert False, "Should have raised TypeError"
except TypeError:
    pass

try:
    test_func(1, 2, 3)  # Too many arguments
    assert False, "Should have raised TypeError"
except TypeError:
    pass

try:
    test_func("invalid", 2)  # Invalid argument type
    assert False, "Should have raised TypeError"
except TypeError:
    pass
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddTypeEdgeCases(t *testing.T) {
	setupTest(t)

	// Test adding non-struct type
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when adding non-struct type")
		}
	}()

	m := MainModule()
	m.AddType(123, nil, "NotAStruct", "")
}

func TestInitFunctionEdgeCases(t *testing.T) {
	setupTest(t)

	type TestStruct struct {
		Value int
	}

	// Test init function with invalid signature
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when using invalid init function")
		}
	}()

	m := MainModule()
	invalidInit := func(x string) string { return x } // Wrong signature
	m.AddType(TestStruct{}, invalidInit, "TestStruct", "")
}

func TestNestedStructRegistration(t *testing.T) {
	setupTest(t)

	type NestedStruct struct {
		Value int
	}

	type ParentStruct struct {
		Nested    NestedStruct
		NestedPtr *NestedStruct
	}

	m := MainModule()
	m.AddType(ParentStruct{}, nil, "ParentStruct", "")

	code := `
parent = ParentStruct()
assert hasattr(parent, "nested")
assert hasattr(parent, "nested_ptr")

# Test nested struct manipulation
parent.nested.value = 42
assert parent.nested.value == 42

parent.nested_ptr = None
assert parent.nested_ptr is None

# Create and assign new nested struct
parent.nested_ptr = NestedStruct()
parent.nested_ptr.value = 100
assert parent.nested_ptr.value == 100
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddTypeWithPointerArg(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type TestStruct struct {
		Value int
	}

	// Test adding type with pointer argument
	typ1 := m.AddType(&TestStruct{}, nil, "TestStruct", "")
	if typ1.Nil() {
		t.Fatal("Failed to create type with pointer argument")
	}

	code := `
obj = TestStruct()
obj.value = 42
assert obj.value == 42
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddTypeDuplicate(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type TestStruct struct {
		Value int
	}

	// First registration
	typ1 := m.AddType(TestStruct{}, nil, "TestStruct", "")
	if typ1.Nil() {
		t.Fatal("Failed to create type on first registration")
	}

	// Second registration should return the same type object
	typ2 := m.AddType(TestStruct{}, nil, "TestStruct", "")
	if typ2.Nil() {
		t.Fatal("Failed to get type on second registration")
	}

	if typ1.cpyObj() != typ2.cpyObj() {
		t.Fatal("Expected same type object on second registration")
	}

	// Both types should work with the same underlying Go type
	code := `
obj1 = TestStruct()
obj1.value = 42
assert obj1.value == 42
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}

	// Also test with pointer argument
	typ3 := m.AddType(&TestStruct{}, nil, "TestStruct3", "")
	if typ3.Nil() {
		t.Fatal("Failed to get type on registration with pointer")
	}

	if typ1.cpyObj() != typ3.cpyObj() {
		t.Fatal("Expected same type object on second registration")
	}
}

func TestStructPointerFieldDictAssignment(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type NestedStruct struct {
		IntVal    int
		StringVal string
	}

	type ParentStruct struct {
		PtrField *NestedStruct
	}

	m.AddType(ParentStruct{}, nil, "ParentStruct", "")

	// Test assigning dict to nil pointer field
	code := `
obj = ParentStruct()
# Initially the pointer should be nil
assert obj.ptr_field is None

# Assign dict to nil pointer field
obj.ptr_field = {"int_val": 42, "string_val": "hello"}
assert obj.ptr_field.int_val == 42
assert obj.ptr_field.string_val == "hello"

# Test invalid dict value type
try:
    obj.ptr_field = {"int_val": "not an int", "string_val": "hello"}
    assert False, "Should have raised TypeError for invalid int_val"
except TypeError:
    pass

# Test completely wrong type
try:
    obj.ptr_field = ["not", "a", "dict"]
    assert False, "Should have raised TypeError for list"
except TypeError:
    pass

# Test nested dict with wrong type
try:
    obj.ptr_field = {"int_val": {"nested": "dict"}, "string_val": "hello"}
    assert False, "Should have raised TypeError for nested dict"
except TypeError:
    pass
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStructFieldDictAssignment(t *testing.T) {
	setupTest(t)
	m := MainModule()

	type NestedStruct struct {
		IntVal    int
		StringVal string
	}

	type ParentStruct struct {
		Field NestedStruct
	}

	m.AddType(ParentStruct{}, nil, "ParentStruct", "")

	// Test assigning dict to struct field
	code := `
obj = ParentStruct()

# Assign valid dict
obj.field = {"int_val": 42, "string_val": "hello"}
assert obj.field.int_val == 42
assert obj.field.string_val == "hello"

# Test invalid value type
try:
    obj.field = {"int_val": "not an int", "string_val": "hello"}
    assert False, "Should have raised TypeError for invalid int_val"
except TypeError:
    pass

# Test completely wrong type
try:
    obj.field = ["not", "a", "dict"]
    assert False, "Should have raised TypeError for list"
except TypeError:
    pass

# Test nested dict with wrong type
try:
    obj.field = {"int_val": {"nested": "dict"}, "string_val": "hello"}
    assert False, "Should have raised TypeError for nested dict"
except TypeError:
    pass

# Test with complex nested structure
obj.field = {
    "int_val": 100,
    "string_val": "test"
}
assert obj.field.int_val == 100
assert obj.field.string_val == "test"
`
	err := RunString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenSig(t *testing.T) {
	type CustomStruct struct {
		Value int
	}

	tests := []struct {
		name        string
		fn          any
		hasRecv     bool
		expectedSig string
	}{
		{
			name:        "function with return",
			fn:          func(x int) string { return "" },
			hasRecv:     false,
			expectedSig: "(arg0, /)",
		},
		{
			name:        "multiple arguments",
			fn:          func(x string, y int) (int, string) { return 0, "" },
			hasRecv:     false,
			expectedSig: "(arg0, arg1, /)",
		},
		{
			name:        "with kwargs",
			fn:          func(name string, kwargs KwArgs) Object { return None() },
			hasRecv:     false,
			expectedSig: "(arg0, /, **kwargs)",
		},
		{
			name:        "method with receiver",
			fn:          func(r *CustomStruct, x int) float64 { return 0 },
			hasRecv:     true,
			expectedSig: "(arg0, /)",
		},
		{
			name:        "no arguments",
			fn:          func() {},
			hasRecv:     false,
			expectedSig: "()",
		},
		{
			name:        "only kwargs",
			fn:          func(kwargs KwArgs) {},
			hasRecv:     false,
			expectedSig: "(**kwargs)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genSig(tt.fn, tt.hasRecv)
			if got != tt.expectedSig {
				t.Errorf("genSig() = %q, want %q", got, tt.expectedSig)
			}
		})
	}
}

func TestCreateFuncNaming(t *testing.T) {
	setupTest(t)

	tests := []struct {
		name         string
		givenName    string
		fn           any
		doc          string
		expectedName string
	}{
		{
			name:         "explicit name",
			givenName:    "my_func",
			fn:           func(x int) int { return x },
			doc:          "test function",
			expectedName: "my_func",
		},
		{
			name:         "empty name uses function name",
			givenName:    "",
			fn:           explicitFunc,
			doc:          "test function",
			expectedName: "explicit_func",
		},
		{
			name:         "camelCase to snake_case",
			givenName:    "camelCaseName",
			fn:           func() {},
			doc:          "test function",
			expectedName: "camel_case_name",
		},
		{
			name:         "package path stripped",
			givenName:    "github.com/user/pkg.MyFunc",
			fn:           func() {},
			doc:          "test function",
			expectedName: "my_func",
		},
	}

	for _, tt := range tests {
		func() {
			f := CreateFunc(tt.givenName, tt.fn, tt.doc)
			if f.Nil() {
				t.Fatal("Failed to create function")
			}

			code := fmt.Sprintf(`
assert "%s" in globals(), "Function %s not found in globals"
assert callable(globals()["%s"]), "Function %s is not callable"
`, tt.expectedName, tt.expectedName, tt.expectedName, tt.expectedName)

			err := RunString(code)
			if err != nil {
				t.Fatalf("Test failed: %v", err)
			}
		}()
	}
}

func TestAddMethodNaming(t *testing.T) {
	setupTest(t)
	m := MainModule()

	tests := []struct {
		name         string
		givenName    string
		fn           any
		doc          string
		expectedName string
	}{
		{
			name:         "explicit name",
			givenName:    "my_method",
			fn:           func(x int) int { return x },
			doc:          "test method",
			expectedName: "my_method",
		},
		{
			name:         "empty name uses function name",
			givenName:    "",
			fn:           explicitFunc,
			doc:          "test method",
			expectedName: "explicit_func",
		},
		{
			name:         "anonymous function gets generated name",
			givenName:    "",
			fn:           func() {},
			doc:          "test method",
			expectedName: "func",
		},
		{
			name:         "camelCase to snake_case",
			givenName:    "myMethodName",
			fn:           func() {},
			doc:          "test method",
			expectedName: "my_method_name",
		},
	}

	for _, tt := range tests {
		func() {
			f := m.AddMethod(tt.givenName, tt.fn, tt.doc)
			if f.Nil() {
				t.Fatal("Failed to create method")
			}
			if !strings.HasPrefix(f.Name(), tt.expectedName) {
				t.Errorf("Expected method name to start with %s, got %s", tt.expectedName, f.Name())
			}
		}()
	}
}
