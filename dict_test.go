package gp

import (
	"testing"
)

func TestDictFromPairs(t *testing.T) {
	setupTest(t)
	// Add panic test case
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("DictFromPairs() with odd number of arguments should panic")
			} else if r != "DictFromPairs requires an even number of arguments" {
				t.Errorf("Expected panic message 'DictFromPairs requires an even number of arguments', got '%v'", r)
			}
		}()

		DictFromPairs("key1", "value1", "key2") // Should panic
	}()

	tests := []struct {
		name     string
		pairs    []any
		wantKeys []any
		wantVals []any
	}{
		{
			name:     "string keys and values",
			pairs:    []any{"key1", "value1", "key2", "value2"},
			wantKeys: []any{"key1", "key2"},
			wantVals: []any{"value1", "value2"},
		},
		{
			name:     "mixed types",
			pairs:    []any{"key1", 42, "key2", 3.14},
			wantKeys: []any{"key1", "key2"},
			wantVals: []any{42, 3.14},
		},
	}

	for _, tt := range tests {
		dict := DictFromPairs(tt.pairs...)

		// Verify each key-value pair
		for i := 0; i < len(tt.wantKeys); i++ {
			key := From(tt.wantKeys[i])
			val := dict.Get(key)
			if !ObjectsAreEqual(val, From(tt.wantVals[i])) {
				t.Errorf("DictFromPairs() got value %v for key %v, want %v",
					val, tt.wantKeys[i], tt.wantVals[i])
			}
		}
	}
}

func TestMakeDict(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name string
		m    map[any]any
	}{
		{
			name: "string map",
			m: map[any]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "mixed types map",
			m: map[any]any{
				"int":    42,
				"float":  3.14,
				"string": "hello",
			},
		},
	}

	for _, tt := range tests {
		dict := MakeDict(tt.m)

		// Verify each key-value pair
		for k, v := range tt.m {
			key := From(k)
			got := dict.Get(key)
			if !ObjectsAreEqual(got, From(v)) {
				t.Errorf("MakeDict() got value %v for key %v, want %v", got, k, v)
			}
		}
	}
}

func TestDictSetGet(t *testing.T) {
	setupTest(t)
	dict := DictFromPairs()

	// Test Set and Get
	key := From("test_key")
	value := From("test_value")
	dict.Set(key, value)

	got := dict.Get(key)
	if !ObjectsAreEqual(got, value) {
		t.Errorf("Dict.Get() got %v, want %v", got, value)
	}
}

func TestDictSetGetString(t *testing.T) {
	setupTest(t)
	dict := DictFromPairs()

	// Test SetString and GetString
	value := From("test_value")
	dict.SetString("test_key", value)

	got := dict.GetString("test_key")
	if !ObjectsAreEqual(got, value) {
		t.Errorf("Dict.GetString() got %v, want %v", got, value)
	}
}

func TestDictDel(t *testing.T) {
	setupTest(t)
	dict := DictFromPairs("test_key", "test_value")
	key := From("test_key")

	// Verify key exists
	got := dict.Get(key)
	if !ObjectsAreEqual(got, From("test_value")) {
		t.Errorf("Before deletion, got %v, want %v", got, "test_value")
	}

	// Delete the key
	dict.Del(key)

	// After deletion, the key should not exist
	if dict.HasKey(key) {
		t.Errorf("After deletion, key %v should not exist", key)
	}
}

func TestDictForEach(t *testing.T) {
	setupTest(t)
	dict := DictFromPairs(
		"key1", "value1",
		"key2", "value2",
		"key3", "value3",
	)

	count := 0
	expectedPairs := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	iter := dict.Iter()
	for iter.HasNext() {
		key, value := iter.Next()
		count++
		k := key.String()
		v := value.String()
		if expectedVal, ok := expectedPairs[k]; !ok || expectedVal != v {
			t.Errorf("ForEach() unexpected pair: %v: %v", k, v)
		}
	}

	if count != len(expectedPairs) {
		t.Errorf("ForEach() visited %d pairs, want %d", count, len(expectedPairs))
	}
}

// Helper function to compare Python objects
func ObjectsAreEqual(obj1, obj2 Object) bool {
	return obj1.String() == obj2.String()
}
