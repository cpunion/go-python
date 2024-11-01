package gp

import (
	"fmt"
	"testing"
)

// TestStruct 包含各种类型的字段用于测试
type TestStruct struct {
	// C兼容的基本类型
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

	// 非C兼容类型
	StringField string
	SliceField  []int
	MapField    map[string]int
	StructField struct{ X int }
}

func (t *TestStruct) TestMethod() int {
	return 42
}

func TestAddType(t *testing.T) {
	Initialize()
	defer Finalize()

	m := MainModule()

	// 测试添加类型
	typ := AddType[TestStruct](m, nil, "TestStruct", "Test struct documentation")
	if typ.Nil() {
		t.Fatal("Failed to create type")
	}

	// 通过Python代码测试类型
	code := `
# 创建实例
obj = TestStruct()

# 测试C兼容类型字段
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

# 测试方法调用
result = obj.test_method()
assert result == 42

# 验证字段值
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

# 测试非C兼容类型字段是否被正确跳过
try:
    obj.string_field = "test"
    assert False, "Should not be able to access string_field"
except AttributeError:
    pass

try:
    obj.slice_field = [1, 2, 3]
    assert False, "Should not be able to access slice_field"
except AttributeError:
    pass

try:
    obj.map_field = {"key": "value"}
    assert False, "Should not be able to access map_field"
except AttributeError:
    pass

try:
    obj.struct_field = None
    assert False, "Should not be able to access struct_field"
except AttributeError:
    pass
`

	err := RunString(code)
	if err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

// 测试带构造函数的类型
type InitTestStruct struct {
	Value int
}

func (i *InitTestStruct) Init(val int) {
	i.Value = val
}

func TestAddTypeWithInit(t *testing.T) {
	Initialize()
	defer Finalize()

	m := MainModule()

	// 测试添加带构造函数的类型
	typ := AddType[InitTestStruct](m, (*InitTestStruct).Init, "InitTestStruct", "Test init struct")
	if typ.Nil() {
		t.Fatal("Failed to create type with init")
	}

	// 通过Python代码测试构造函数
	code := `
# 测试构造函数
obj = InitTestStruct(42)
print(dir(obj))
assert obj.value == 42

# 测试无参数调用时的错误处理
try:
    obj2 = InitTestStruct()
    assert False, "Should fail without arguments"
except TypeError:
    print("========= ")

# 测试参数类型错误的处理
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

// RunString executes Python code string and returns error if any
func RunString(code string) error {
	// Get __main__ module dict for executing code
	main := MainModule()
	dict := main.Dict()

	// Run the code string
	codeObj := CompileString(code, "<string>", FileInput)
	if codeObj.Nil() {
		return fmt.Errorf("failed to compile code")
	}

	ret := EvalCode(codeObj, dict, dict)
	if ret.Nil() {
		if err := FetchError(); err != nil {
			return err
		}
		return fmt.Errorf("failed to execute code")
	}
	return nil
}
