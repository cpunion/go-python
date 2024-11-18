package install

import (
	"path/filepath"
)

const (
	// DepsDir is the directory for all dependencies
	DepsDir = ".deps"
	// PythonDir is the directory name for Python installation
	PythonDir = "python"
)

// GetPythonRoot returns the Python installation root path relative to project path
func GetPythonRoot(projectPath string) string {
	return filepath.Join(projectPath, DepsDir, PythonDir)
}

// GetPythonBinDir returns the Python binary directory path relative to project path
func GetPythonBinDir(projectPath string) string {
	return filepath.Join(GetPythonRoot(projectPath), "bin")
}

// GetPythonLibDir returns the Python library directory path relative to project path
func GetPythonLibDir(projectPath string) string {
	return filepath.Join(GetPythonRoot(projectPath), "lib")
}

// GetPythonPkgConfigDir returns the pkg-config directory path relative to project path
func GetPythonPkgConfigDir(projectPath string) string {
	return filepath.Join(GetPythonLibDir(projectPath), "pkgconfig")
}
