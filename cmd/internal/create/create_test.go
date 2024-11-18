package create

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temporary directory for testing
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "gopy-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

// cleanupTestDir removes the temporary test directory
func cleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("failed to cleanup test dir: %v", err)
	}
}

func TestProject(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(dir string) error
		wantErr bool
	}{
		{
			name: "create new project in empty directory",
			setup: func(dir string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "create project with existing directory",
			setup: func(dir string) error {
				return os.MkdirAll(dir, 0755)
			},
			wantErr: false,
		},
		{
			name: "create project with existing files",
			setup: func(dir string) error {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
				// Create a go.mod file
				return os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := setupTestDir(t)
			defer cleanupTestDir(t, testDir)

			// Setup test case
			if err := tt.setup(testDir); err != nil {
				t.Fatalf("test setup failed: %v", err)
			}

			// Run Project function
			err := Project(testDir, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Project() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify project structure
			expectedFiles := []string{
				"go.mod",
				"main.go",
				".gitignore",
			}

			for _, file := range expectedFiles {
				path := filepath.Join(testDir, file)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("expected file %s does not exist", file)
				}
			}
		})
	}
}

func TestPromptOverwrite(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantOverwrite    bool
		wantOverwriteAll bool
	}{
		{
			name:             "answer yes",
			input:            "y\n",
			wantOverwrite:    true,
			wantOverwriteAll: false,
		},
		{
			name:             "answer no",
			input:            "n\n",
			wantOverwrite:    false,
			wantOverwriteAll: false,
		},
		{
			name:             "answer all",
			input:            "a\n",
			wantOverwrite:    true,
			wantOverwriteAll: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file to simulate stdin
			tmpfile, err := os.CreateTemp("", "test-input")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write test input to temp file
			if _, err := tmpfile.Write([]byte(tt.input)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Redirect stdin to our temp file
			oldStdin := os.Stdin
			f, err := os.Open(tmpfile.Name())
			if err != nil {
				t.Fatal(err)
			}
			os.Stdin = f
			defer func() {
				os.Stdin = oldStdin
				f.Close()
			}()

			// Test promptOverwrite
			gotOverwrite, gotOverwriteAll := promptOverwrite("test.txt")
			if gotOverwrite != tt.wantOverwrite {
				t.Errorf("promptOverwrite() overwrite = %v, want %v", gotOverwrite, tt.wantOverwrite)
			}
			if gotOverwriteAll != tt.wantOverwriteAll {
				t.Errorf("promptOverwrite() overwriteAll = %v, want %v", gotOverwriteAll, tt.wantOverwriteAll)
			}
		})
	}
}
