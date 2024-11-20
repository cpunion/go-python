package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cpunion/go-python/internal/env"
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
		if runtime.GOOS == "windows" {
			os.Setenv("USERPROFILE", tmpDir)
		} else {
			os.Setenv("HOME", tmpDir)
		}

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
		if runtime.GOOS == "windows" {
			os.Setenv("USERPROFILE", "/nonexistent/path")
		} else {
			os.Setenv("HOME", "/nonexistent/path")
		}

		_, err := getCacheDir()
		if err == nil {
			t.Error("getCacheDir() error = nil, want error for invalid home directory")
		}
	})
}

func TestUpdatePkgConfig(t *testing.T) {
	t.Run("freethreaded pkg-config files", func(t *testing.T) {
		// Create temporary directory structure
		tmpDir := t.TempDir()
		pkgConfigDir := env.GetPythonPkgConfigDir(tmpDir)
		if err := os.MkdirAll(pkgConfigDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create test .pc files with freethreaded content
		testFiles := map[string]string{
			"python-3.13t.pc": `prefix=/install
libdir=${prefix}/lib
includedir=${prefix}/include

Name: Python
Description: Python library
Version: 3.13
Libs: -L${libdir} -lpython3t
Cflags: -I${includedir}`,
			"python-3.13t-embed.pc": `prefix=/install
libdir=${prefix}/lib
includedir=${prefix}/include

Name: Python
Description: Embed Python into an application
Version: 3.13
Libs: -L${libdir} -lpython313t
Cflags: -I${includedir}`,
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
		expectedFiles := map[string]struct {
			shouldExist bool
			libName     string
		}{
			// Freethreaded versions
			"python-3.13t.pc":       {true, "-lpython3t"},
			"python3t.pc":           {true, "-lpython3t"},
			"python-3.13t-embed.pc": {true, "-lpython313t"},
			"python3t-embed.pc":     {true, "-lpython313t"},
			// Non-t versions (same content as freethreaded)
			"python-3.13.pc":       {true, "-lpython3t"},
			"python3.pc":           {true, "-lpython3t"},
			"python-3.13-embed.pc": {true, "-lpython313t"},
			"python3-embed.pc":     {true, "-lpython313t"},
		}

		absPath, _ := filepath.Abs(filepath.Join(tmpDir, ".deps/python"))
		for filename, expected := range expectedFiles {
			path := filepath.Join(pkgConfigDir, filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if expected.shouldExist {
					t.Errorf("Expected file %s was not created", filename)
				}
				continue
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", filename, err)
				continue
			}

			// Check prefix
			expectedPrefix := fmt.Sprintf("prefix=%s", absPath)
			if !strings.Contains(string(content), expectedPrefix) {
				t.Errorf("File %s does not contain expected prefix %s", filename, expectedPrefix)
			}

			// Check library name
			if !strings.Contains(string(content), expected.libName) {
				t.Errorf("File %s does not contain expected library name %s", filename, expected.libName)
			}
		}
	})

	t.Run("non-freethreaded pkg-config files", func(t *testing.T) {
		// Create temporary directory structure
		tmpDir := t.TempDir()
		pkgConfigDir := env.GetPythonPkgConfigDir(tmpDir)
		if err := os.MkdirAll(pkgConfigDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create test .pc files with non-freethreaded content
		testFiles := map[string]string{
			"python-3.13.pc": `prefix=/install
libdir=${prefix}/lib
includedir=${prefix}/include

Name: Python
Description: Python library
Version: 3.13
Libs: -L${libdir} -lpython3
Cflags: -I${includedir}`,
			"python-3.13-embed.pc": `prefix=/install
libdir=${prefix}/lib
includedir=${prefix}/include

Name: Python
Description: Embed Python into an application
Version: 3.13
Libs: -L${libdir} -lpython313
Cflags: -I${includedir}`,
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
		expectedFiles := map[string]struct {
			shouldExist bool
			libName     string
		}{
			"python-3.13.pc":       {true, "-lpython3"},
			"python3.pc":           {true, "-lpython3"},
			"python-3.13-embed.pc": {true, "-lpython313"},
			"python3-embed.pc":     {true, "-lpython313"},
		}

		absPath, _ := filepath.Abs(filepath.Join(tmpDir, ".deps/python"))
		for filename, expected := range expectedFiles {
			path := filepath.Join(pkgConfigDir, filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if expected.shouldExist {
					t.Errorf("Expected file %s was not created", filename)
				}
				continue
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", filename, err)
				continue
			}

			// Check prefix
			expectedPrefix := fmt.Sprintf("prefix=%s", absPath)
			if !strings.Contains(string(content), expectedPrefix) {
				t.Errorf("File %s does not contain expected prefix %s", filename, expectedPrefix)
			}

			// Check library name
			if !strings.Contains(string(content), expected.libName) {
				t.Errorf("File %s does not contain expected library name %s", filename, expected.libName)
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
