package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// depsDir is the directory for all dependencies
	depsDir = ".deps"
	// pyDir is the directory name for Python installation
	pyDir = "python"
	// goDir is the directory name for Go installation
	goDir = "go"
	// mingwDir is the directory name for Mingw installation
	mingwDir  = "mingw"
	mingwRoot = mingwDir + "/mingw64"

	tinyPkgConfigDir = "tiny-pkg-config"
)

func GetDepsDir(projectPath string) string {
	return filepath.Join(projectPath, depsDir)
}

func GetGoDir(projectPath string) string {
	return filepath.Join(GetDepsDir(projectPath), goDir)
}

// GetPythonRoot returns the Python installation root path relative to project path
func GetPythonRoot(projectPath string) string {
	return filepath.Join(projectPath, depsDir, pyDir)
}

// GetPythonBinDir returns the Python binary directory path relative to project path
func GetPythonBinDir(projectPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(GetPythonRoot(projectPath))
	}
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
	return filepath.Join(projectPath, depsDir, goDir)
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

func GetMingwDir(projectPath string) string {
	return filepath.Join(projectPath, depsDir, mingwDir)
}

func GetMingwRoot(projectPath string) string {
	return filepath.Join(projectPath, depsDir, mingwRoot)
}

func GetTinyPkgConfigDir(projectPath string) string {
	return filepath.Join(projectPath, depsDir, tinyPkgConfigDir)
}

func GetEnvConfigPath(projectPath string) string {
	return filepath.Join(GetDepsDir(projectPath), "env.txt")
}

func SetBuildEnv(projectPath string) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		panic(err)
	}
	path := os.Getenv("PATH")
	path = GetGoBinDir(absPath) + pathSeparator() + path
	path = GetPythonBinDir(absPath) + pathSeparator() + path
	if runtime.GOOS == "windows" {
		path = GetMingwRoot(absPath) + pathSeparator() + path
		path = GetTinyPkgConfigDir(absPath) + pathSeparator() + path
	}
	os.Setenv("PATH", path)
	os.Setenv("GOPATH", GetGoPath(absPath))
	os.Setenv("GOROOT", GetGoRoot(absPath))
	os.Setenv("GOCACHE", GetGoCacheDir(absPath))
	os.Setenv("PKG_CONFIG_PATH", GetPythonPkgConfigDir(absPath))
	os.Setenv("CGO_ENABLED", "1")
}

func pathSeparator() string {
	if runtime.GOOS == "windows" {
		return ";"
	}
	return ":"
}

// WriteEnvFile writes environment variables to .deps/env.txt
func WriteEnvFile(projectPath, pythonHome, pythonPath string) error {
	// Prepare environment variables
	envVars := []string{
		fmt.Sprintf("PYTHONPATH=%s", strings.TrimSpace(pythonPath)),
		fmt.Sprintf("PYTHONHOME=%s", pythonHome),
		fmt.Sprintf("PATH=%s", GetPythonBinDir(projectPath)),
	}

	// Write to env.txt
	envFile := GetEnvConfigPath(projectPath)
	if err := os.WriteFile(envFile, []byte(strings.Join(envVars, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write env file: %v", err)
	}

	return nil
}

// ReadEnvFile loads environment variables from .python/env.txt in the given directory
func ReadEnvFile(projectDir string) (map[string]string, error) {
	envFile := GetEnvConfigPath(projectDir)
	content, err := os.ReadFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file %s: %v", envFile, err)
	}
	envs := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envs[parts[0]] = parts[1]
		}
	}
	return envs, nil
}

func GeneratePythonEnv(pythonHome, pythonPath string) map[string]string {
	path := os.Getenv("PATH")
	if runtime.GOOS == "windows" {
		path = filepath.Join(pythonHome) + ";" + path
	} else {
		path = filepath.Join(pythonHome, "bin") + ":" + path
	}
	return map[string]string{
		"PYTHONHOME": pythonHome,
		"PYTHONPATH": pythonPath,
		"PATH":       path,
	}
}

func ReadEnv(projectDir string) (map[string]string, error) {
	envs, err := ReadEnvFile(projectDir)
	if err != nil {
		return nil, err
	}
	pythonHome, ok := envs["PYTHONHOME"]
	if !ok {
		return nil, fmt.Errorf("PYTHONHOME is not set in env.txt")
	}
	pythonPath, ok := envs["PYTHONPATH"]
	if !ok {
		return nil, fmt.Errorf("PYTHONPATH is not set in env.txt")
	}
	for k, v := range GeneratePythonEnv(pythonHome, pythonPath) {
		envs[k] = v
	}
	return envs, nil
}
