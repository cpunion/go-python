package gp

import (
	"testing"
)

func TestMakeStr(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name      string
		input     string
		expected  string
		length    int
		byteCount int
	}{
		{
			name:      "empty string",
			input:     "",
			expected:  "",
			length:    0,
			byteCount: 0,
		},
		{
			name:      "ascii string",
			input:     "hello",
			expected:  "hello",
			length:    5,
			byteCount: 5, // ASCII character each takes 1 byte
		},
		{
			name:      "unicode string",
			input:     "你好世界",
			expected:  "你好世界",
			length:    4,
			byteCount: 12, // Chinese character each takes 3 bytes
		},
		{
			name:      "mixed string",
			input:     "hello世界",
			expected:  "hello世界",
			length:    7,
			byteCount: 11, // 5 ASCII characters (5 bytes) + 2 Chinese characters (6 bytes)
		},
		{
			name:      "special unicode",
			input:     "π∑€",
			expected:  "π∑€",
			length:    3,
			byteCount: 8, // π(2bytes) + ∑(3bytes) + €(3bytes) = 8bytes
		},
	}

	for _, tt := range tests {
		pyStr := MakeStr(tt.input)

		// Test String() method
		if got := pyStr.String(); got != tt.expected {
			t.Errorf("MakeStr(%q).String() = %q, want %q", tt.input, got, tt.expected)
		}

		// Test Len() method
		if got := pyStr.Len(); got != tt.length {
			t.Errorf("MakeStr(%q).Len() = %d, want %d", tt.input, got, tt.length)
		}

		// Test ByteLen() method
		if got := pyStr.ByteLen(); got != tt.byteCount {
			t.Errorf("MakeStr(%q).ByteLen() = %d, want %d", tt.input, got, tt.byteCount)
		}
	}
}

func TestStrEncode(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name     string
		input    string
		encoding string
	}{
		{
			name:     "utf-8 encoding",
			input:    "hello世界",
			encoding: "utf-8",
		},
		{
			name:     "ascii encoding",
			input:    "hello",
			encoding: "ascii",
		},
	}

	for _, tt := range tests {
		pyStr := MakeStr(tt.input)
		encoded := pyStr.Encode(tt.encoding)
		decoded := encoded.Decode(tt.encoding)

		if got := decoded.String(); got != tt.input {
			t.Errorf("String encode/decode roundtrip failed: got %q, want %q", got, tt.input)
		}
	}
}
