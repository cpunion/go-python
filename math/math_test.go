package math

import (
	"runtime"
	"testing"

	gp "github.com/cpunion/go-python"
)

func TestSqrt(t *testing.T) {
	runtime.LockOSThread()
	// Initialize Python
	gp.Initialize()
	defer gp.Finalize()

	tests := []struct {
		input    float64
		expected float64
	}{
		{16.0, 4.0},
		{25.0, 5.0},
		{0.0, 0.0},
		{100.0, 10.0},
	}

	for _, test := range tests {
		input := gp.MakeFloat(test.input)
		result := Sqrt(input)

		if result.Float64() != test.expected {
			t.Errorf("Sqrt(%f) = %f; want %f", test.input, result.Float64(), test.expected)
		}
	}
}
