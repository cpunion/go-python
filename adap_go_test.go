package gp

import (
	"testing"
)

func TestAllocCStr(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"ascii string", "hello", "hello"},
		{"unicode string", "hello 世界", "hello 世界"},
	}

	for _, tt := range tests {
		cstr := AllocCStr(tt.input)
		got := GoString(cstr)
		if got != tt.want {
			t.Errorf("AllocCStr() = %v, want %v", got, tt.want)
		}
	}
}

func TestGoStringN(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name  string
		input string
		n     int
		want  string
	}{
		{"empty string", "", 0, ""},
		{"partial string", "hello", 3, "hel"},
		{"full string", "hello", 5, "hello"},
		{"unicode partial", "hello 世界", 6, "hello "},
		{"unicode full", "hello 世界", 12, "hello 世界"},
	}

	for _, tt := range tests {
		cstr := AllocCStr(tt.input)
		got := GoStringN(cstr, tt.n)
		if got != tt.want {
			t.Errorf("GoStringN() = %v, want %v", got, tt.want)
		}
	}
}

func TestAllocCStrDontFree(t *testing.T) {
	setupTest(t)
	input := "test string"
	cstr := AllocCStrDontFree(input)
	got := GoString(cstr)
	if got != input {
		t.Errorf("AllocCStrDontFree() = %v, want %v", got, input)
	}
}
