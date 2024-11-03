package gp

import (
	"reflect"
	"testing"
)

func TestObjectCreation(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

func TestObjectAttributes(t *testing.T) {
	// Test attributes using Python's built-in object type
	builtins := ImportModule("builtins")
	obj := builtins.AttrFunc("object").Call()

	t.Run("GetAttr", func(t *testing.T) {
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
	})

	t.Run("SetAttr", func(t *testing.T) {
		// Create a custom class to test attribute setting
		pyCode := `
class TestClass:
    pass
`
		locals := MakeDict(nil)
		globals := MakeDict(nil)
		globals.Set(MakeStr("__builtins__"), builtins.Object)

		code := CompileString(pyCode, "<string>", FileInput)

		EvalCode(code, globals, locals).AsModule()
		testClass := locals.Get(MakeStr("TestClass")).AsFunc()
		instance := testClass.Call()

		// Now we can set attributes
		instance.SetAttr("new_attr", "test_value")
		value := instance.Attr("new_attr")
		if value.AsStr().String() != "test_value" {
			t.Error("SetAttr failed to set new attribute")
		}
	})
}

func TestDictOperations(t *testing.T) {
	// Test dictionary operations
	pyDict := MakeDict(nil)
	pyDict.Set(MakeStr("key1"), From(42))
	pyDict.Set(MakeStr("key2"), From("value"))

	t.Run("GetItem", func(t *testing.T) {
		value := pyDict.Get(MakeStr("key1"))
		if value.AsLong().Int64() != 42 {
			t.Error("Failed to get dictionary item")
		}
	})

	t.Run("SetItem", func(t *testing.T) {
		pyDict.Set(MakeStr("key3"), From("new_value"))
		value := pyDict.Get(MakeStr("key3"))
		if value.AsStr().String() != "new_value" {
			t.Error("Failed to set dictionary item")
		}
	})
}

func TestObjectConversion(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("StructConversion", func(t *testing.T) {
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
	})

	t.Run("SliceConversion", func(t *testing.T) {
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
	})
}

func TestObjectString(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			obj := From(tt.input)
			str := obj.String()
			if str != tt.expected {
				t.Errorf("String() = %v, want %v", str, tt.expected)
			}
		})
	}
}
