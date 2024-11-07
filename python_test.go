package gp

import (
	"runtime"
	"sync"
	"testing"
)

var (
	testMutex sync.Mutex
)

func setupTest(t *testing.T) {
	testMutex.Lock()
	Initialize()
	// TODO: Remove this once we solve random segfaults
	t.Cleanup(func() {
		runtime.GC()
		Finalize()
		testMutex.Unlock()
	})
}

func TestRunString(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "valid python code",
			code:    "x = 1 + 1",
			wantErr: false,
		},
		{
			name:    "invalid python code",
			code:    "x = ",
			wantErr: true,
		},
		{
			name:    "syntax error",
			code:    "for i in range(10) print(i)", // missing :
			wantErr: true,
		},
	}

	for _, tt := range tests {
		err := RunString(tt.code)
		if (err != nil) != tt.wantErr {
			t.Errorf("RunString() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestCompileString(t *testing.T) {
	setupTest(t)
	tests := []struct {
		name     string
		code     string
		filename string
		start    InputType
		wantNil  bool
	}{
		{
			name:     "compile expression",
			code:     "1 + 1",
			filename: "<string>",
			start:    EvalInput,
			wantNil:  false,
		},
		{
			name:     "compile invalid code",
			code:     "x =",
			filename: "<string>",
			start:    EvalInput,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		obj, _ := CompileString(tt.code, tt.filename, tt.start)
		if obj.Nil() != tt.wantNil {
			t.Errorf("CompileString() returned nil = %v, want %v", obj.Nil(), tt.wantNil)
		}
	}
}

func TestNone(t *testing.T) {
	setupTest(t)
	none := None()
	if none.Nil() {
		t.Error("None() returned nil object")
	}
}

func TestNil(t *testing.T) {
	setupTest(t)
	nil_ := Nil()
	if !nil_.Nil() {
		t.Error("Nil() did not return nil object")
	}
}

func TestMainModule(t *testing.T) {
	setupTest(t)
	main := MainModule()
	if main.Nil() {
		t.Error("MainModule() returned nil")
	}
}

func TestWith(t *testing.T) {
	setupTest(t)
	// First create a simple Python context manager class
	code := `
class TestContextManager:
    def __init__(self):
        self.entered = False
        self.exited = False
    
    def __enter__(self):
        self.entered = True
        return self
    
    def __exit__(self, *args):
        self.exited = True
        return None
`
	if err := RunString(code); err != nil {
		t.Fatalf("Failed to create test context manager: %v", err)
	}

	// Get the context manager class and create an instance
	main := MainModule()
	cmClass := main.AttrFunc("TestContextManager")
	if cmClass.Nil() {
		t.Fatal("Failed to get TestContextManager class")
	}

	cm := cmClass.Call()
	if cm.Nil() {
		t.Fatal("Failed to create context manager instance")
	}

	// Test the With function
	called := false
	With(cm, func(obj Object) {
		called = true

		// Check that __enter__ was called
		entered := obj.AttrBool("entered")
		if entered.Nil() {
			t.Error("Could not get entered attribute")
		}
		if !entered.Bool() {
			t.Error("__enter__ was not called")
		}
	})

	// Verify the callback was called
	if !called {
		t.Error("With callback was not called")
	}

	// Check that __exit__ was called
	exited := cm.AttrBool("exited")
	if exited.Nil() {
		t.Error("Could not get exited attribute")
	}
	if !exited.Bool() {
		t.Error("__exit__ was not called")
	}
}
