package gp

import (
	"testing"
)

func TestModuleImport(t *testing.T) {
	Initialize()
	defer Finalize()

	// Test importing a built-in module
	mathMod := ImportModule("math")
	if mathMod.Nil() {
		t.Fatal("Failed to import math module")
	}

	// Test getting module dictionary
	modDict := mathMod.Dict()
	if modDict.Nil() {
		t.Fatal("Failed to get module dictionary")
	}

	// Verify math module has expected attributes
	if !modDict.Has("pi") {
		t.Error("Math module doesn't contain 'pi' constant")
	}
}

func TestGetModule(t *testing.T) {
	Initialize()
	defer Finalize()

	// First import the module
	sysModule := ImportModule("sys")
	if sysModule.Nil() {
		t.Fatal("Failed to import sys module")
	}

	// Then try to get it
	gotModule := GetModule("sys")
	if gotModule.Nil() {
		t.Fatal("Failed to get sys module")
	}

	// Both should refer to the same module
	if !sysModule.Equals(gotModule) {
		t.Error("GetModule returned different module instance than ImportModule")
	}
}

func TestCreateModule(t *testing.T) {
	Initialize()
	defer Finalize()

	// Create a new module
	modName := "test_module"
	mod := CreateModule(modName)
	if mod.Nil() {
		t.Fatal("Failed to create new module")
	}

	// Add an object to the module
	value := From(42)
	err := mod.AddObject("test_value", value)
	if err != 0 {
		t.Fatal("Failed to add object to module")
	}

	// Verify the object was added
	modDict := mod.Dict()
	if !modDict.Has("test_value") {
		t.Error("Module doesn't contain added value")
	}

	// Verify the value is correct
	gotValue := modDict.Get(From("test_value"))
	if !gotValue.Equals(value) {
		t.Error("Retrieved value doesn't match added value")
	}
}

func TestGetModuleDict(t *testing.T) {
	Initialize()
	defer Finalize()

	// Get the module dictionary
	moduleDict := GetModuleDict()
	if moduleDict.Nil() {
		t.Fatal("Failed to get module dictionary")
	}

	// Import a module
	mathMod := ImportModule("math")
	if mathMod.Nil() {
		t.Fatal("Failed to import math module")
	}

	// Verify the module is in the module dictionary
	if !moduleDict.Has("math") {
		t.Error("Module dictionary doesn't contain imported module")
	}
}
