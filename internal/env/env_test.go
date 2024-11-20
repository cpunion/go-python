package env

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	t.Run("valid env file", func(t *testing.T) {
		// Create temporary directory structure
		projectDir := t.TempDir()
		pythonDir := GetPythonRoot(projectDir)
		if err := os.MkdirAll(pythonDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create test env.txt file
		envContent := map[string]string{
			"PKG_CONFIG_PATH": "/test/lib/pkgconfig",
			"PYTHONPATH":      "/test/lib/python3.9",
			"PYTHONHOME":      "/test",
		}
		lines := []string{}
		for key, value := range envContent {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
		}
		envFile := GetEnvConfigPath(projectDir)
		if err := os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			t.Fatal(err)
		}

		// Test loading the env file
		got, err := ReadEnvFile(projectDir)
		if err != nil {
			t.Errorf("LoadEnvFile() error = %v, want nil", err)
			return
		}

		if !reflect.DeepEqual(got, envContent) {
			t.Errorf("LoadEnvFile() = %v, want %v", got, envContent)
		}
	})

	t.Run("missing env file", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := ReadEnvFile(tmpDir)
		if err == nil {
			t.Error("LoadEnvFile() error = nil, want error for missing env file")
		}
	})
}

func TestWriteEnvFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Run("write env file", func(t *testing.T) {
		// Create temporary directory structure
		projectDir := t.TempDir()
		pythonDir := GetPythonRoot(projectDir)
		binDir := GetPythonBinDir(projectDir)
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create mock Python executable
		var pythonPath string
		if runtime.GOOS == "windows" {
			pythonPath = filepath.Join(binDir, "python.exe")
			pythonScript := `@echo off
echo /mock/path1;/mock/path2
`
			if err := os.WriteFile(pythonPath, []byte(pythonScript), 0644); err != nil {
				t.Fatal(err)
			}
		} else {
			pythonPath = filepath.Join(binDir, "python")
			pythonScript := `#!/bin/sh
echo "/mock/path1:/mock/path2"
`
			if err := os.WriteFile(pythonPath, []byte(pythonScript), 0755); err != nil {
				t.Fatal(err)
			}
		}

		// Test writing env file
		if err := WriteEnvFile(projectDir); err != nil {
			t.Errorf("writeEnvFile() error = %v, want nil", err)
			return
		}

		// Verify the env file was created
		envFile := GetEnvConfigPath(projectDir)
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			t.Error("writeEnvFile() did not create env.txt")
			return
		}

		// Read and verify content
		content, err := os.ReadFile(envFile)
		if err != nil {
			t.Errorf("Failed to read env.txt: %v", err)
			return
		}

		// Get expected path separator
		pathSep := ":"
		if runtime.GOOS == "windows" {
			pathSep = ";"
		}

		// Verify the content contains expected environment variables
		envContent := string(content)
		expectedVars := []string{
			fmt.Sprintf("PKG_CONFIG_PATH=%s", filepath.Join(pythonDir, "lib", "pkgconfig")),
			fmt.Sprintf("PYTHONPATH=/mock/path1%s/mock/path2", pathSep),
			fmt.Sprintf("PYTHONHOME=%s", pythonDir),
		}
		fmt.Printf("envContent:\n%v\n", envContent)
		for _, v := range expectedVars {
			if !strings.Contains(envContent, v) {
				t.Errorf("env.txt missing expected variable %s", v)
			}
		}
	})

	t.Run("missing python executable", func(t *testing.T) {
		tmpDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(tmpDir, ".deps/python"), 0755); err != nil {
			t.Fatal(err)
		}

		err := WriteEnvFile(tmpDir)
		if err == nil {
			t.Error("writeEnvFile() error = nil, want error for missing python executable")
		}
	})
}
