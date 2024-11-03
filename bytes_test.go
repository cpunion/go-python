package gp

import (
	"bytes"
	"testing"
)

func TestBytesCreation(t *testing.T) {
	// Test BytesFromStr
	b1 := BytesFromStr("hello")
	if string(b1.Bytes()) != "hello" {
		t.Errorf("BytesFromStr: expected 'hello', got '%s'", string(b1.Bytes()))
	}

	// Test MakeBytes
	data := []byte("world")
	b2 := MakeBytes(data)
	if !bytes.Equal(b2.Bytes(), data) {
		t.Errorf("MakeBytes: expected '%v', got '%v'", data, b2.Bytes())
	}
}

func TestBytesDecode(t *testing.T) {
	// Test UTF-8 decode
	b := BytesFromStr("你好")
	if !bytes.Equal(b.Bytes(), []byte("你好")) {
		t.Errorf("BytesFromStr: expected '你好', got '%s'", string(b.Bytes()))
	}
	s := b.Decode("utf-8")
	if s.String() != "你好" {
		t.Errorf("Decode: expected '你好', got '%s'", s.String())
	}

	// Test ASCII decode
	b2 := BytesFromStr("hello")
	s2 := b2.Decode("ascii")
	if s2.String() != "hello" {
		t.Errorf("Decode: expected 'hello', got '%s'", s2.String())
	}
}

func TestBytesConversion(t *testing.T) {
	original := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" in hex
	b := MakeBytes(original)

	// Test conversion back to []byte
	result := b.Bytes()
	if !bytes.Equal(result, original) {
		t.Errorf("Bytes conversion: expected %v, got %v", original, result)
	}
}
