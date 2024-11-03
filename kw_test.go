package gp

import (
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name    string
		args    []any
		wantTup Tuple
		wantKw  KwArgs
	}{
		{
			name:    "empty args",
			args:    []any{},
			wantTup: MakeTuple(),
			wantKw:  nil,
		},
		{
			name:    "only positional args",
			args:    []any{1, "two", 3.0},
			wantTup: MakeTuple(1, "two", 3.0),
			wantKw:  nil,
		},
		{
			name:    "with kwargs",
			args:    []any{1, "two", KwArgs{"a": 1, "b": "test"}},
			wantTup: MakeTuple(1, "two"),
			wantKw:  KwArgs{"a": 1, "b": "test"},
		},
		{
			name:    "only kwargs",
			args:    []any{KwArgs{"x": 10, "y": 20}},
			wantTup: MakeTuple(),
			wantKw:  KwArgs{"x": 10, "y": 20},
		},
	}

	for _, tt := range tests {
		gotTup, gotKw := splitArgs(tt.args...)

		if !reflect.DeepEqual(gotTup, tt.wantTup) {
			t.Errorf("splitArgs() tuple = %v, want %v", gotTup, tt.wantTup)
		}

		if !reflect.DeepEqual(gotKw, tt.wantKw) {
			t.Errorf("splitArgs() kwargs = %v, want %v", gotKw, tt.wantKw)
		}
	}
}

func TestKwArgs(t *testing.T) {
	setupTest(t)
	kw := KwArgs{
		"name": "test",
		"age":  42,
	}

	// Test type assertion
	if _, ok := interface{}(kw).(KwArgs); !ok {
		t.Error("KwArgs failed type assertion")
	}

	// Test map operations
	kw["new"] = "value"
	if v, ok := kw["new"]; !ok || v != "value" {
		t.Error("KwArgs map operations failed")
	}
}
