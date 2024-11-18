package install

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestGetPythonURL(t *testing.T) {
	tests := []struct {
		name         string
		arch         string
		os           string
		freeThreaded bool

		debug   bool
		want    string
		wantErr bool
	}{
		{
			name:         "darwin-arm64-freethreaded-debug",
			arch:         "arm64",
			os:           "darwin",
			freeThreaded: true,
			debug:        true,
			want:         "cpython-3.13.0+20241016-aarch64-apple-darwin-freethreaded+debug-full.tar.zst",
		},
		{
			name:         "darwin-amd64-freethreaded-pgo",
			arch:         "amd64",
			os:           "darwin",
			freeThreaded: true,
			debug:        false,
			want:         "cpython-3.13.0+20241016-x86_64-apple-darwin-freethreaded+pgo-full.tar.zst",
		},
		{
			name:         "darwin-amd64-debug",
			arch:         "amd64",
			os:           "darwin",
			freeThreaded: false,
			debug:        true,
			want:         "cpython-3.13.0+20241016-x86_64-apple-darwin-debug-full.tar.zst",
		},
		{
			name:         "darwin-amd64-pgo",
			arch:         "amd64",
			os:           "darwin",
			freeThreaded: false,
			debug:        false,
			want:         "cpython-3.13.0+20241016-x86_64-apple-darwin-pgo-full.tar.zst",
		},
		{
			name:         "linux-amd64-freethreaded-debug",
			arch:         "amd64",
			os:           "linux",
			freeThreaded: true,
			debug:        true,
			want:         "cpython-3.13.0+20241016-x86_64-unknown-linux-gnu-freethreaded+debug-full.tar.zst",
		},
		{
			name:         "windows-amd64-freethreaded-pgo",
			arch:         "amd64",
			os:           "windows",
			freeThreaded: true,
			debug:        false,
			want:         "cpython-3.13.0+20241016-x86_64-pc-windows-msvc-shared-freethreaded+pgo-full.tar.zst",
		},
		{
			name:         "windows-386-freethreaded-pgo",
			arch:         "386",
			os:           "windows",
			freeThreaded: true,
			debug:        false,
			want:         "cpython-3.13.0+20241016-i686-pc-windows-msvc-shared-freethreaded+pgo-full.tar.zst",
		},
		{
			name:         "unsupported-arch",
			arch:         "mips",
			os:           "linux",
			freeThreaded: false,
			debug:        false,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "unsupported-os",
			arch:         "amd64",
			os:           "freebsd",
			freeThreaded: false,
			debug:        false,
			want:         "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPythonURL("3.13.0", "20241016", tt.arch, tt.os, tt.freeThreaded, tt.debug)

			if tt.wantErr {
				if got != "" {
					t.Errorf("getPythonURL() = %v, want empty string for error case", got)
				}
				return
			}

			if got == "" {
				t.Errorf("getPythonURL() returned empty string, want %v", tt.want)
				return
			}

			// Extract filename from URL
			parts := strings.Split(got, "/")
			filename := parts[len(parts)-1]

			if filename != tt.want {
				t.Errorf("getPythonURL() = %v, want %v", filename, tt.want)
			}
		})
	}
}

func TestGetCacheDir(t *testing.T) {
	// Save original home dir
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	t.Run("valid home directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv("HOME", tmpDir)

		got, err := getCacheDir()
		if err != nil {
			t.Errorf("getCacheDir() error = %v, want nil", err)
			return
		}

		want := filepath.Join(tmpDir, ".gopy", "cache")
		if got != want {
			t.Errorf("getCacheDir() = %v, want %v", got, want)
		}

		// Verify directory was created
		if _, err := os.Stat(got); os.IsNotExist(err) {
			t.Errorf("getCacheDir() did not create cache directory")
		}
	})

	t.Run("invalid home directory", func(t *testing.T) {
		// Set HOME to a non-existent directory
		os.Setenv("HOME", "/nonexistent/path")

		_, err := getCacheDir()
		if err == nil {
			t.Error("getCacheDir() error = nil, want error for invalid home directory")
		}
	})
}

func TestLoadEnvFile(t *testing.T) {
	t.Run("valid env file", func(t *testing.T) {
		// Create temporary directory structure
		tmpDir := t.TempDir()
		pythonDir := GetPythonRoot(tmpDir)
		if err := os.MkdirAll(pythonDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create test env.txt file
		envContent := []string{
			"PKG_CONFIG_PATH=/test/lib/pkgconfig",
			"PYTHONPATH=/test/lib/python3.9",
			"PYTHONHOME=/test",
		}
		envFile := filepath.Join(pythonDir, "env.txt")
		if err := os.WriteFile(envFile, []byte(strings.Join(envContent, "\n")), 0644); err != nil {
			t.Fatal(err)
		}

		// Test loading the env file
		got, err := LoadEnvFile(tmpDir)
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
		_, err := LoadEnvFile(tmpDir)
		if err == nil {
			t.Error("LoadEnvFile() error = nil, want error for missing env file")
		}
	})
}

func TestUpdatePkgConfig(t *testing.T) {
	t.Run("valid pkg-config files", func(t *testing.T) {
		// Create temporary directory structure
		tmpDir := t.TempDir()
		pkgConfigDir := GetPythonPkgConfigDir(tmpDir)
		if err := os.MkdirAll(pkgConfigDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create test .pc files
		testFiles := map[string]string{
			"python-3.13t.pc":      "prefix=/install\nlibdir=${prefix}/lib\n",
			"python-3.13-embed.pc": "prefix=/install\nlibdir=${prefix}/lib\n",
		}

		for filename, content := range testFiles {
			if err := os.WriteFile(filepath.Join(pkgConfigDir, filename), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}

		// Test updating pkg-config files
		if err := updatePkgConfig(tmpDir); err != nil {
			t.Errorf("updatePkgConfig() error = %v, want nil", err)
			return
		}

		// Verify the generated files
		expectedFiles := []string{
			"python-3.13t.pc",
			"python-3.13.pc",
			"python3t.pc",
			"python3.pc",
			"python-3.13-embed.pc",
			"python3-embed.pc",
		}

		for _, filename := range expectedFiles {
			path := filepath.Join(pkgConfigDir, filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", filename)
				continue
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", filename, err)
				continue
			}

			absPath, _ := filepath.Abs(filepath.Join(tmpDir, ".deps/python"))
			expectedPrefix := fmt.Sprintf("prefix=%s", absPath)
			if !strings.Contains(string(content), expectedPrefix) {
				t.Errorf("File %s does not contain expected prefix %s", filename, expectedPrefix)
			}
		}
	})

	t.Run("missing pkgconfig directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := updatePkgConfig(tmpDir)
		if err == nil {
			t.Error("updatePkgConfig() error = nil, want error for missing pkgconfig directory")
		}
	})
}

func TestWriteEnvFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Run("write env file", func(t *testing.T) {
		// Create temporary directory structure
		tmpDir := t.TempDir()
		pythonDir := GetPythonRoot(tmpDir)
		binDir := GetPythonBinDir(tmpDir)
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
		if err := writeEnvFile(tmpDir); err != nil {
			t.Errorf("writeEnvFile() error = %v, want nil", err)
			return
		}

		// Verify the env file was created
		envFile := filepath.Join(pythonDir, "env.txt")
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

		err := writeEnvFile(tmpDir)
		if err == nil {
			t.Error("writeEnvFile() error = nil, want error for missing python executable")
		}
	})
}
