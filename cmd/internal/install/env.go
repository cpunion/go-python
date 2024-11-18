package install

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	// DepsDir is the directory for all dependencies
	DepsDir = ".deps"
	// PythonDir is the directory name for Python installation
	PythonDir = "python"
	// GoDir is the directory name for Go installation
	GoDir = "go"
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

// GetGoRoot returns the Go installation root path relative to project path
func GetGoRoot(projectPath string) string {
	return filepath.Join(projectPath, DepsDir, GoDir)
}

// GetGoPath returns the Go path relative to project path
func GetGoPath(projectPath string) string {
	return filepath.Join(GetGoRoot(projectPath), "packages")
}

// GetGoBinDir returns the Go binary directory path relative to project path
func GetGoBinDir(projectPath string) string {
	return filepath.Join(GetGoRoot(projectPath), "bin")
}

// GetGoCacheDir returns the Go cache directory path relative to project path
func GetGoCacheDir(projectPath string) string {
	return filepath.Join(GetGoRoot(projectPath), "go-build")
}

func SetEnv(projectPath string) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		panic(err)
	}
	os.Setenv("PATH", GetGoBinDir(absPath)+pathSeparator()+os.Getenv("PATH"))
	os.Setenv("GOPATH", GetGoPath(absPath))
	os.Setenv("GOROOT", GetGoRoot(absPath))
	os.Setenv("GOCACHE", GetGoCacheDir(absPath))
}

func pathSeparator() string {
	if runtime.GOOS == "windows" {
		return ";"
	}
	return ":"
}
