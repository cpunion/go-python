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
assert obj.struct_field["x"] == 100

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
assert obj.struct_field["x"] == 100
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
