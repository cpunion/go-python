package gp

import (
	"testing"
)

func TestFloat(t *testing.T) {
	setupTest(t)
	func() {
		// Test creating float and converting back
		f := MakeFloat(3.14159)

		// Test Float64 conversion
		if got := f.Float64(); got != 3.14159 {
			t.Errorf("Float64() = %v, want %v", got, 3.14159)
		}

		// Test Float32 conversion
		if got := f.Float32(); float64(got) != float64(float32(3.14159)) {
			t.Errorf("Float32() = %v, want %v", got, float32(3.14159))
		}
	}()

	func() {
		// Test integer float
		intFloat := MakeFloat(5.0)

		if !intFloat.IsInteger().Bool() {
			t.Errorf("IsInteger() for 5.0 = false, want true")
		}

		// Test non-integer float
		fracFloat := MakeFloat(5.5)

		if fracFloat.IsInteger().Bool() {
			t.Errorf("IsInteger() for 5.5 = true, want false")
		}
	}()

	func() {
		// Test zero
		zero := MakeFloat(0.0)

		if got := zero.Float64(); got != 0.0 {
			t.Errorf("Float64() = %v, want 0.0", got)
		}

		// Test very large number
		large := MakeFloat(1e308)

		if got := large.Float64(); got != 1e308 {
			t.Errorf("Float64() = %v, want 1e308", got)
		}
	}()
}
