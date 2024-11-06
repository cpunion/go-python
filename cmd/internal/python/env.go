package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Env represents a Python environment
type Env struct {
	Root string // Root directory of the Python installation
}

// New creates a new Python environment instance
func New(projectPath string) *Env {
	return &Env{
		Root: filepath.Join(projectPath, ".python"),
	}
}

// Python returns the path to the Python executable
func (e *Env) Python() (string, error) {
	if runtime.GOOS == "windows" {
		pythonPath := filepath.Join(e.Root, "bin", "python3.exe")
		if _, err := os.Stat(pythonPath); err == nil {
			return pythonPath, nil
		}
		pythonPath = filepath.Join(e.Root, "bin", "python.exe")
		if _, err := os.Stat(pythonPath); err == nil {
			return pythonPath, nil
		}
	} else {
		pythonPath := filepath.Join(e.Root, "bin", "python3")
		if _, err := os.Stat(pythonPath); err == nil {
			return pythonPath, nil
		}
		pythonPath = filepath.Join(e.Root, "bin", "python")
		if _, err := os.Stat(pythonPath); err == nil {
			return pythonPath, nil
		}
	}
	return "", fmt.Errorf("python executable not found in %s", e.Root)
}

// Pip returns the path to the pip executable
func (e *Env) Pip() (string, error) {
	if runtime.GOOS == "windows" {
		pipPath := filepath.Join(e.Root, "bin", "pip3.exe")
		if _, err := os.Stat(pipPath); err == nil {
			return pipPath, nil
		}
		pipPath = filepath.Join(e.Root, "bin", "pip.exe")
		if _, err := os.Stat(pipPath); err == nil {
			return pipPath, nil
		}
	} else {
		pipPath := filepath.Join(e.Root, "bin", "pip3")
		if _, err := os.Stat(pipPath); err == nil {
			return pipPath, nil
		}
		pipPath = filepath.Join(e.Root, "bin", "pip")
		if _, err := os.Stat(pipPath); err == nil {
			return pipPath, nil
		}
	}
	return "", fmt.Errorf("pip executable not found in %s", e.Root)
}

// RunPip executes pip with the given arguments
func (e *Env) RunPip(args ...string) error {
	pipPath, err := e.Pip()
	if err != nil {
		return err
	}

	cmd := exec.Command(pipPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
