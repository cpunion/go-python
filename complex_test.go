package gp

import (
	"testing"
)

func TestComplex(t *testing.T) {
	tests := []struct {
		name     string
		input    complex128
		wantReal float64
		wantImag float64
	}{
		{
			name:     "zero complex",
			input:    complex(0, 0),
			wantReal: 0,
			wantImag: 0,
		},
		{
			name:     "positive real and imaginary",
			input:    complex(3.14, 2.718),
			wantReal: 3.14,
			wantImag: 2.718,
		},
		{
			name:     "negative real and imaginary",
			input:    complex(-1.5, -2.5),
			wantReal: -1.5,
			wantImag: -2.5,
		},
		{
			name:     "mixed signs",
			input:    complex(-1.23, 4.56),
			wantReal: -1.23,
			wantImag: 4.56,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MakeComplex(tt.input)

			// Test Real() method
			if got := c.Real(); got != tt.wantReal {
				t.Errorf("Complex.Real() = %v, want %v", got, tt.wantReal)
			}

			// Test Imag() method
			if got := c.Imag(); got != tt.wantImag {
				t.Errorf("Complex.Imag() = %v, want %v", got, tt.wantImag)
			}

			// Test Complex128() method
			if got := c.Complex128(); got != tt.input {
				t.Errorf("Complex.Complex128() = %v, want %v", got, tt.input)
			}
		})
	}
}

func TestComplexZeroValue(t *testing.T) {
	// Create a proper zero complex number instead of using zero-value struct
	c := MakeComplex(complex(0, 0))

	// Test that zero complex behaves correctly
	if got := c.Real(); got != 0 {
		t.Errorf("Zero Complex.Real() = %v, want 0", got)
	}
	if got := c.Imag(); got != 0 {
		t.Errorf("Zero Complex.Imag() = %v, want 0", got)
	}
	if got := c.Complex128(); got != 0 {
		t.Errorf("Zero Complex.Complex128() = %v, want 0", got)
	}
}

func TestComplexNilHandling(t *testing.T) {
	var c Complex // zero-value struct with nil pointer
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil pointer access, but got none")
		}
	}()

	// This should panic
	_ = c.Real()
}
