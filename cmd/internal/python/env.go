package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

// Env represents a Python environment
type Env struct {
	Root string // Root directory of the Python installation
}

// New creates a new Python environment instance
func New(pythonHome string) *Env {
	return &Env{
		Root: pythonHome,
	}
}

// Python returns the path to the Python executable
func (e *Env) Python() (string, error) {
	binDir := e.Root
	if runtime.GOOS != "windows" {
		binDir = filepath.Join(e.Root, "bin")
	}
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return "", fmt.Errorf("failed to read bin directory: %v", err)
	}

	// Single pattern to match all variants, prioritizing 't' versions
	var pattern *regexp.Regexp
	if runtime.GOOS == "windows" {
		pattern = regexp.MustCompile(`^python3?[\d.]*t?(?:\.exe)?$`)
	} else {
		pattern = regexp.MustCompile(`^python3?[\d.]*t?$`)
	}

	for _, entry := range entries {
		if !entry.IsDir() && pattern.MatchString(entry.Name()) {
			return filepath.Join(binDir, entry.Name()), nil
		}
	}

	return "", fmt.Errorf("python executable not found in %s", e.Root)
}

// RunPip executes pip with the given arguments
func (e *Env) RunPip(args ...string) error {
	return e.RunPython(append([]string{"-m", "pip"}, args...)...)
}

// RunPython executes python with the given arguments
func (e *Env) RunPython(args ...string) error {
	pythonPath, err := e.Python()
	if err != nil {
		return err
	}

	cmd := exec.Command(pythonPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
