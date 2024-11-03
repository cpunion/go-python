package gp

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Initialize()
	code := m.Run()
	Finalize()
	os.Exit(code)
}
