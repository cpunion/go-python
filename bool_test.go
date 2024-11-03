package gp

import (
	"testing"
)

func TestBool(t *testing.T) {
	// Test MakeBool
	b1 := MakeBool(true)
	if !b1.Bool() {
		t.Error("MakeBool(true) should return true")
	}

	b2 := MakeBool(false)
	if b2.Bool() {
		t.Error("MakeBool(false) should return false")
	}

	// Test True and False
	if !True().Bool() {
		t.Error("True() should return true")
	}

	if False().Bool() {
		t.Error("False() should return false")
	}

	// Test Not method
	if True().Not().Bool() {
		t.Error("True().Not() should return false")
	}

	if !False().Not().Bool() {
		t.Error("False().Not() should return true")
	}
}
