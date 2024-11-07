package gp

import (
	"testing"
)

func TestTupleCreation(t *testing.T) {
	setupTest(t)
	// Test empty tuple
	empty := MakeTupleWithLen(0)
	if empty.Len() != 0 {
		t.Errorf("Expected empty tuple length 0, got %d", empty.Len())
	}

	// Test tuple with values
	tuple := MakeTuple(42, "hello", 3.14)
	if tuple.Len() != 3 {
		t.Errorf("Expected tuple length 3, got %d", tuple.Len())
	}
}

func TestTupleGetSet(t *testing.T) {
	setupTest(t)
	tuple := MakeTupleWithLen(2)

	// Test setting and getting values
	tuple.Set(0, From(123))
	tuple.Set(1, From("test"))

	if val := tuple.Get(0).AsLong().Int64(); val != 123 {
		t.Errorf("Expected 123, got %d", val)
	}
	if val := tuple.Get(1).AsStr().String(); val != "test" {
		t.Errorf("Expected 'test', got %s", val)
	}
}

func TestTupleSlice(t *testing.T) {
	setupTest(t)
	tuple := MakeTuple(1, 2, 3, 4, 5)

	// Test slicing
	slice := tuple.Slice(1, 4)
	if slice.Len() != 3 {
		t.Errorf("Expected slice length 3, got %d", slice.Len())
	}

	expected := []int64{2, 3, 4}
	for i := 0; i < slice.Len(); i++ {
		if val := slice.Get(i).AsLong().Int64(); val != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], val)
		}
	}
}

func TestTupleParseArgs(t *testing.T) {
	setupTest(t)
	tuple := MakeTuple(42, "hello", 3.14, true)

	var (
		intVal   int
		strVal   string
		floatVal float64
		boolVal  bool
		extraVal int // This shouldn't get set
	)

	// Test successful parsing
	success := tuple.ParseArgs(&intVal, &strVal, &floatVal, &boolVal)
	if !success {
		t.Error("ParseArgs failed unexpectedly")
	}

	if intVal != 42 {
		t.Errorf("Expected int 42, got %d", intVal)
	}
	if strVal != "hello" {
		t.Errorf("Expected string 'hello', got %s", strVal)
	}
	if floatVal != 3.14 {
		t.Errorf("Expected float 3.14, got %f", floatVal)
	}
	if !boolVal {
		t.Errorf("Expected bool true, got false")
	}

	// Test parsing with too many arguments
	success = tuple.ParseArgs(&intVal, &strVal, &floatVal, &boolVal, &extraVal)
	if success {
		t.Error("ParseArgs should have failed with too many arguments")
	}

	// Test parsing with invalid type
	var invalidPtr *testing.T
	success = tuple.ParseArgs(&invalidPtr)
	if success {
		t.Error("ParseArgs should have failed with invalid type")
	}
}

func TestTupleParseArgsTypes(t *testing.T) {
	setupTest(t)
	// Test all supported numeric types
	tuple := MakeTuple(42, 42, 42, 42, 42, 42, 42, 42, 42, 42)

	var (
		intVal    int
		int8Val   int8
		int16Val  int16
		int32Val  int32
		int64Val  int64
		uintVal   uint
		uint8Val  uint8
		uint16Val uint16
		uint32Val uint32
		uint64Val uint64
	)

	success := tuple.ParseArgs(
		&intVal, &int8Val, &int16Val, &int32Val, &int64Val,
		&uintVal, &uint8Val, &uint16Val, &uint32Val, &uint64Val,
	)

	if !success {
		t.Error("ParseArgs failed for numeric types")
	}

	// Test floating point types
	floatTuple := MakeTuple(3.14, 3.14)
	var float32Val float32
	var float64Val float64

	success = floatTuple.ParseArgs(&float32Val, &float64Val)
	if !success {
		t.Error("ParseArgs failed for floating point types")
	}

	// Test complex types
	complexTuple := MakeTuple(complex(1, 2), complex(3, 4))
	var complex64Val complex64
	var complex128Val complex128

	success = complexTuple.ParseArgs(&complex64Val, &complex128Val)
	if !success {
		t.Error("ParseArgs failed for complex types")
	}

	// Test string and bytes
	strTuple := MakeTuple("hello")
	var strVal string
	var bytesVal []byte
	var objVal Object
	var pyObj *cPyObject

	success = strTuple.ParseArgs(&strVal)
	if !success || strVal != "hello" {
		t.Error("ParseArgs failed for string type")
	}

	success = strTuple.ParseArgs(&bytesVal)
	if !success || string(bytesVal) != "hello" {
		t.Error("ParseArgs failed for bytes type")
	}

	success = strTuple.ParseArgs(&objVal)
	if !success || !objVal.IsStr() {
		t.Error("ParseArgs failed for object type")
	}

	success = strTuple.ParseArgs(&pyObj)
	if !success || pyObj == nil {
		t.Error("ParseArgs failed for PyObject type")
	}
	str := FromPy(pyObj)
	if !str.IsStr() || str.String() != "hello" {
		t.Error("FromPy returned non-string object")
	}
}
